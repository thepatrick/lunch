'use strict';

const ExtractTextPlugin = require("extract-text-webpack-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const SriPlugin = require('webpack-subresource-integrity');
const webpackUglifyJsPlugin = require('webpack-uglify-js-plugin');
const webpack = require('webpack');

module.exports = {
  devtool: 'sourcemap',
  plugins: [
    new ExtractTextPlugin({ filename: 'css/[name].css', disable: false, allChunks: true }),
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery",
      "window.jQuery": "jquery",
      Tether: "tether",
      "window.Tether": "tether",
      Alert: "exports-loader?Alert!bootstrap/js/dist/alert",
      Button: "exports-loader?Button!bootstrap/js/dist/button",
      Carousel: "exports-loader?Carousel!bootstrap/js/dist/carousel",
      Collapse: "exports-loader?Collapse!bootstrap/js/dist/collapse",
      Dropdown: "exports-loader?Dropdown!bootstrap/js/dist/dropdown",
      Modal: "exports-loader?Modal!bootstrap/js/dist/modal",
      Popover: "exports-loader?Popover!bootstrap/js/dist/popover",
      Scrollspy: "exports-loader?Scrollspy!bootstrap/js/dist/scrollspy",
      Tab: "exports-loader?Tab!bootstrap/js/dist/tab",
      Tooltip: "exports-loader?Tooltip!bootstrap/js/dist/tooltip",
      Util: "exports-loader?Util!bootstrap/js/dist/util",
    }),
    new HtmlWebpackPlugin({
      title: 'Lunch Bot',
    }),
    new SriPlugin({
      hashFuncNames: ['sha256', 'sha384'],
      enabled: process.env.NODE_ENV === 'production',
    }),
    new webpack.DefinePlugin({
      "process.env.NODE_ENV": JSON.stringify(process.env.NODE_ENV || 'development'),
    })
  ],
  resolve: {
    extensions: ['.js', '.jsx'],
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
      },
      // the url-loader uses DataUrls.
      // the file-loader emits files.
      {
        test: /\.(woff|woff2)$/,
        loader: 'url-loader?limit=10000&mimetype=application/font-woff',
      },
      { test: /\.ttf$/, loader: 'file-loader' },
      { test: /\.eot$/, loader: 'file-loader' },
      { test: /\.svg$/, loader: 'file-loader' },
      {
        test: /\/bootstrap\/js\//,
        loader: 'imports-loader?jQuery=jquery',
      },
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract({
          fallbackLoader: "style-loader",
          loader: "css-loader"
        })
      }
    ],
  },
  entry: {
    app: ['./client/index.jsx'],
  },
  output: {
    path: path.resolve(__dirname, 'static', 'manage'),
    filename: 'app.js',
    publicPath: '/static/manage/', 
    crossOriginLoading: 'anonymous',
  },
};

if (process.env.NODE_ENV === 'production') {
  module.exports.plugins = module.exports.plugins.concat(
      new webpackUglifyJsPlugin({
        enabled: process.env.NODE_ENV === 'production',
        cacheFolder: path.resolve(__dirname, 'static/cached_uglify/'),
        debug: process.env.NODE_ENV !== 'production',
        minimize: process.env.NODE_ENV === 'production',
        sourceMap: process.env.NODE_ENV !== 'production',
        output: {
          comments: process.env.NODE_ENV !== 'production'
        },
        compressor: {
          warnings: process.env.NODE_ENV === 'production'
        }
      })
  );
}