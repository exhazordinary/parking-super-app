import path from 'path';
import { fileURLToPath } from 'url';
import * as Repack from '@callstack/repack';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default (env) => {
  const { mode = 'development', context = __dirname, entry = './src/index.ts', platform = process.env.PLATFORM || 'ios', devServer = undefined } = env;
  return {
    mode, devtool: mode === 'development' ? 'source-map' : false, context, entry,
    resolve: { ...Repack.getResolveOptions(platform) },
    output: { clean: true, hashFunction: 'xxhash64', path: path.resolve(__dirname, 'build/generated', platform), filename: 'index.bundle', publicPath: Repack.getPublicPath({ platform, devServer }) },
    module: { rules: [
      { test: /\.[jt]sx?$/, include: [/node_modules(.*[/\\])+react-native/], use: 'babel-loader' },
      { test: /\.[jt]sx?$/, exclude: /node_modules/, use: { loader: 'babel-loader', options: { presets: ['module:@react-native/babel-preset'] } } },
      { test: Repack.getAssetExtensionsRegExp(Repack.ASSET_EXTENSIONS.filter((ext) => ext !== 'svg')), use: { loader: '@callstack/repack/assets-loader', options: { platform, devServerEnabled: Boolean(devServer) } } },
    ]},
    plugins: [
      new Repack.RepackPlugin({ context, mode, platform, devServer }),
      new Repack.plugins.ModuleFederationPlugin({
        name: 'notification_module', exposes: { './NotificationNavigator': './src/NotificationNavigator' },
        shared: { react: { singleton: true }, 'react-native': { singleton: true }, 'react-native-paper': { singleton: true }, '@react-navigation/native': { singleton: true }, '@tanstack/react-query': { singleton: true }, '@parking/ui': { singleton: true }, '@parking/api': { singleton: true }, '@parking/navigation': { singleton: true } },
      }),
    ],
  };
};
