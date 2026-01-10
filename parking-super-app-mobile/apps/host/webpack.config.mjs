import path from 'path';
import { fileURLToPath } from 'url';
import TerserPlugin from 'terser-webpack-plugin';
import * as Repack from '@callstack/repack';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Webpack configuration for the host app.
 * This is the Module Federation host that loads remote modules.
 */
export default (env) => {
  const {
    mode = 'development',
    context = __dirname,
    entry = './index.js',
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
    entry: [
      ...Repack.getInitializationEntries(reactNativePath, {
        hmr: devServer && devServer.hmr,
      }),
      entry,
    ],
    resolve: {
      ...Repack.getResolveOptions(platform),
      alias: {
        '@': path.resolve(dirname, 'src'),
      },
    },
    output: {
      clean: true,
      hashFunction: 'xxhash64',
      path: path.resolve(dirname, 'build/generated', platform),
      filename: 'index.bundle',
      chunkFilename: '[name].chunk.bundle',
      publicPath: Repack.getPublicPath({ platform, devServer }),
    },
    optimization: {
      minimize,
      minimizer: [
        new TerserPlugin({
          test: /\.(js)?bundle(\?.*)?$/i,
          extractComments: false,
          terserOptions: {
            format: {
              comments: false,
            },
          },
        }),
      ],
      chunkIds: 'named',
    },
    module: {
      rules: [
        {
          test: /\.[jt]sx?$/,
          include: [
            /node_modules(.*[/\\])+react-native/,
            /node_modules(.*[/\\])+@react-native/,
            /node_modules(.*[/\\])+@react-navigation/,
            /node_modules(.*[/\\])+react-native-reanimated/,
            /node_modules(.*[/\\])+react-native-gesture-handler/,
            /node_modules(.*[/\\])+react-native-screens/,
            /node_modules(.*[/\\])+react-native-safe-area-context/,
            /node_modules(.*[/\\])+react-native-paper/,
            /node_modules(.*[/\\])+react-native-vector-icons/,
            /node_modules(.*[/\\])+@callstack[/\\]repack/,
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
              plugins: ['react-native-reanimated/plugin'],
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
              scalableAssetExtensions: Repack.SCALABLE_ASSETS,
            },
          },
        },
        {
          test: /\.svg$/,
          use: [
            {
              loader: '@svgr/webpack',
              options: {
                native: true,
              },
            },
          ],
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
        name: 'host',
        shared: {
          react: {
            singleton: true,
            eager: true,
            requiredVersion: '18.2.0',
          },
          'react-native': {
            singleton: true,
            eager: true,
            requiredVersion: '0.73.6',
          },
          'react-native-paper': {
            singleton: true,
            eager: true,
          },
          'react-native-safe-area-context': {
            singleton: true,
            eager: true,
          },
          'react-native-screens': {
            singleton: true,
            eager: true,
          },
          'react-native-gesture-handler': {
            singleton: true,
            eager: true,
          },
          'react-native-reanimated': {
            singleton: true,
            eager: true,
          },
          '@react-navigation/native': {
            singleton: true,
            eager: true,
          },
          '@react-navigation/native-stack': {
            singleton: true,
            eager: true,
          },
          '@react-navigation/bottom-tabs': {
            singleton: true,
            eager: true,
          },
          '@tanstack/react-query': {
            singleton: true,
            eager: true,
          },
          zustand: {
            singleton: true,
            eager: true,
          },
          '@parking/ui': {
            singleton: true,
            eager: true,
          },
          '@parking/api': {
            singleton: true,
            eager: true,
          },
          '@parking/auth': {
            singleton: true,
            eager: true,
          },
          '@parking/navigation': {
            singleton: true,
            eager: true,
          },
        },
      }),
    ],
  };
};
