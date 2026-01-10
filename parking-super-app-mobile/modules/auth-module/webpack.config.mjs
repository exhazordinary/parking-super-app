import path from 'path';
import { fileURLToPath } from 'url';
import * as Repack from '@callstack/repack';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Webpack configuration for auth-module (Module Federation remote)
 */
export default (env) => {
  const {
    mode = 'development',
    context = __dirname,
    entry = './src/index.ts',
    platform = process.env.PLATFORM || 'ios',
    minimize = mode === 'production',
    devServer = undefined,
    bundleFilename = undefined,
    sourceMapFilename = undefined,
    assetsPath = undefined,
    reactNativePath = require.resolve('react-native'),
  } = env;

  const dirname = context;

  return {
    mode,
    devtool: mode === 'development' ? 'source-map' : false,
    context: dirname,
    entry,
    resolve: {
      ...Repack.getResolveOptions(platform),
    },
    output: {
      clean: true,
      hashFunction: 'xxhash64',
      path: path.resolve(dirname, 'build/generated', platform),
      filename: 'index.bundle',
      chunkFilename: '[name].chunk.bundle',
      publicPath: Repack.getPublicPath({ platform, devServer }),
    },
    module: {
      rules: [
        {
          test: /\.[jt]sx?$/,
          include: [
            /node_modules(.*[/\\])+react-native/,
            /node_modules(.*[/\\])+@react-native/,
            /node_modules(.*[/\\])+@react-navigation/,
          ],
          use: 'babel-loader',
        },
        {
          test: /\.[jt]sx?$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['module:@react-native/babel-preset'],
            },
          },
        },
        {
          test: Repack.getAssetExtensionsRegExp(
            Repack.ASSET_EXTENSIONS.filter((ext) => ext !== 'svg')
          ),
          use: {
            loader: '@callstack/repack/assets-loader',
            options: {
              platform,
              devServerEnabled: Boolean(devServer),
            },
          },
        },
      ],
    },
    plugins: [
      new Repack.RepackPlugin({
        context: dirname,
        mode,
        platform,
        devServer,
        output: {
          bundleFilename,
          sourceMapFilename,
          assetsPath,
        },
      }),
      new Repack.plugins.ModuleFederationPlugin({
        name: 'auth_module',
        exposes: {
          './AuthNavigator': './src/AuthNavigator',
        },
        shared: {
          react: {
            singleton: true,
            eager: false,
            requiredVersion: '18.2.0',
          },
          'react-native': {
            singleton: true,
            eager: false,
            requiredVersion: '0.73.6',
          },
          'react-native-paper': {
            singleton: true,
            eager: false,
          },
          '@react-navigation/native': {
            singleton: true,
            eager: false,
          },
          '@react-navigation/native-stack': {
            singleton: true,
            eager: false,
          },
          '@tanstack/react-query': {
            singleton: true,
            eager: false,
          },
          zustand: {
            singleton: true,
            eager: false,
          },
          '@parking/ui': {
            singleton: true,
            eager: false,
          },
          '@parking/api': {
            singleton: true,
            eager: false,
          },
          '@parking/auth': {
            singleton: true,
            eager: false,
          },
          '@parking/navigation': {
            singleton: true,
            eager: false,
          },
        },
      }),
    ],
  };
};
