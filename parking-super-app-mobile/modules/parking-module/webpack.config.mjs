import path from 'path';
import { fileURLToPath } from 'url';
import * as Repack from '@callstack/repack';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export default (env) => {
  const { mode = 'development', context = __dirname, entry = './src/index.ts', platform = process.env.PLATFORM || 'ios', devServer = undefined, bundleFilename = undefined, sourceMapFilename = undefined, assetsPath = undefined } = env;

  return {
    mode,
    devtool: mode === 'development' ? 'source-map' : false,
    context,
    entry,
    resolve: { ...Repack.getResolveOptions(platform) },
    output: { clean: true, hashFunction: 'xxhash64', path: path.resolve(__dirname, 'build/generated', platform), filename: 'index.bundle', chunkFilename: '[name].chunk.bundle', publicPath: Repack.getPublicPath({ platform, devServer }) },
    module: {
      rules: [
        { test: /\.[jt]sx?$/, include: [/node_modules(.*[/\\])+react-native/, /node_modules(.*[/\\])+@react-native/], use: 'babel-loader' },
        { test: /\.[jt]sx?$/, exclude: /node_modules/, use: { loader: 'babel-loader', options: { presets: ['module:@react-native/babel-preset'] } } },
        { test: Repack.getAssetExtensionsRegExp(Repack.ASSET_EXTENSIONS.filter((ext) => ext !== 'svg')), use: { loader: '@callstack/repack/assets-loader', options: { platform, devServerEnabled: Boolean(devServer) } } },
      ],
    },
    plugins: [
      new Repack.RepackPlugin({ context, mode, platform, devServer, output: { bundleFilename, sourceMapFilename, assetsPath } }),
      new Repack.plugins.ModuleFederationPlugin({
        name: 'parking_module',
        exposes: { './ParkingNavigator': './src/ParkingNavigator' },
        shared: { react: { singleton: true, eager: false }, 'react-native': { singleton: true, eager: false }, 'react-native-paper': { singleton: true, eager: false }, '@react-navigation/native': { singleton: true, eager: false }, '@react-navigation/native-stack': { singleton: true, eager: false }, '@tanstack/react-query': { singleton: true, eager: false }, '@parking/ui': { singleton: true, eager: false }, '@parking/api': { singleton: true, eager: false }, '@parking/auth': { singleton: true, eager: false }, '@parking/navigation': { singleton: true, eager: false } },
      }),
    ],
  };
};
