const path = require('path');
const HtmlWebpackPlugin = require("html-webpack-plugin");
const HtmlWebpackHarddiskPlugin = require('html-webpack-harddisk-plugin');
const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');
const CopyPlugin = require("copy-webpack-plugin");

module.exports = (argv) => {
  const isDevelopment = argv.mode === 'development';
  const isDevServer = process.env?.BUILD !== '1';

  return {
    target: 'web',
    entry: './src/index.tsx',
    output: {
      filename: isDevelopment ? '[name].js' : '[name].[contenthash].js',
      path: path.resolve(__dirname, 'build'),
      clean: true,
      assetModuleFilename: '[path]/[name][ext]'
    },
    optimization: {
      // https://webpack.js.org/guides/caching/#extracting-boilerplate
      runtimeChunk: 'single',
      splitChunks: {
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            chunks: 'all',
          },
        },
      },
    },
    module: {
      rules: [
        {
          test: /\.(js|tsx?)$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: [
                '@babel/preset-env',
                '@babel/preset-react',
                '@babel/preset-typescript'
              ],
              plugins: [
                isDevServer && require.resolve('react-refresh/babel'),
                // Required for async/await:
                ['@babel/transform-runtime', { regenerator: true }],
              ].filter(Boolean),
            }
          }
        },
        {
          test: /\.svg$/,
          exclude: /node_modules/,
          type: 'asset/resource',
          use: 'svgo-loader',
        },
        {
          test: /\.(jpeg|jpg|png|woff|woff2|eot|ttf)(\?.*$|$)/,
          exclude: /node_modules/,
          type: 'asset',
        },
        {
          // https://www.npmjs.com/package/sass-loader
          test: /\.s[ac]ss$/i,
          exclude: /node_modules/,
          use: [
            "style-loader",
            "css-loader",
            "sass-loader",
          ],
        },
      ]
    },
    plugins: [
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, 'public/index.html'),
        publicPath: '/',
        alwaysWriteToDisk: true, // See below
      }),
      // Forces the index.html generation, even when using webpack-dev-server
      // https://github.com/jantimon/html-webpack-harddisk-plugin
      isDevServer && new HtmlWebpackHarddiskPlugin(),
      // Enables fast-refresh for true React HMR:
      // https://github.com/pmmmwh/react-refresh-webpack-plugin
      isDevServer && new ReactRefreshWebpackPlugin(),
      new CopyPlugin({
        patterns: [
          {
            from: 'public',
            globOptions: {
              ignore: ['**/index.{html,ejs}'],
            }
          },
        ],
      }),
    ].filter(Boolean),
    resolve: {
      extensions: ['.ts', '.tsx', '.js'],
      alias: {
        '@app': path.resolve(__dirname, 'src'),
        '@assets': path.resolve(__dirname, 'assets'),
        '@images': path.resolve(__dirname, 'assets/images'),
      }
    },
    devtool: isDevServer ? 'eval-cheap-module-source-map' : 'source-map',
    devServer: {
      // https://webpack.js.org/configuration/dev-server/#devserverhot
      hot: true,
      // https://webpack.js.org/configuration/dev-server/#devserversetupmiddlewares
      historyApiFallback: true,
      compress: true,
      port: 'auto',
    },
  }
}
