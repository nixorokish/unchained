import { config as getEnv, parse } from 'dotenv'
import { existsSync, readFileSync } from 'fs'
import { join } from 'path'
import * as k8s from '@pulumi/kubernetes'
import { deployApi } from '@shapeshiftoss/common-pulumi'
import { deployIndexer } from '@shapeshiftoss/blockbook-pulumi'
import { getConfig } from './config'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Outputs = Record<string, any>

//https://www.pulumi.com/docs/intro/languages/javascript/#entrypoint
export = async (): Promise<Outputs> => {
  const { kubeconfig, config, namespace } = await getConfig()
  const { cluster } = config

  const name = 'unchained'
  const asset = config.network !== 'mainnet' ? `dogecoin-${config.network}` : 'dogecoin'
  const outputs: Outputs = {}

  let provider: k8s.Provider
  if (config.isLocal) {
    provider = new k8s.Provider('kube-provider', { cluster, context: cluster })
  } else {
    provider = new k8s.Provider('kube-provider', { kubeconfig })
  }

  if (existsSync('../.env')) {
    getEnv({ path: join(__dirname, '../.env') })
  } else if (config.isLocal) {
    throw new Error(
      'you must run `cp sample.env .env` from the dogecoin coinstack directory and fill out any empty values.'
    )
  }

  const missingKeys: Array<string> = []
  const stringData = Object.keys(parse(readFileSync('../sample.env'))).reduce((prev, key) => {
    const value = process.env[key]

    if (!value) {
      missingKeys.push(key)
      return prev
    }

    return { ...prev, [key]: value }
  }, {})

  if (missingKeys.length) {
    throw new Error(`Missing the following required environment variables: ${missingKeys.join(', ')}`)
  }

  new k8s.core.v1.Secret(asset, { metadata: { name: asset, namespace }, stringData }, { provider })

  await deployIndexer(name, asset, provider, namespace, config)
  await deployApi(name, asset, provider, namespace, config)

  return outputs
}
