const path = require('path');
const webpack = require('webpack');

/*
 * SplitChunksPlugin is enabled by default and replaced
 * deprecated CommonsChunkPlugin. It automatically identifies modules which
 * should be splitted of chunk by heuristics using module duplication count and
 * module category (i. e. node_modules). And splits the chunks…
 *
 * It is safe to remove "splitChunks" from the generated configuration
 * and was added as an educational example.
 *
 * https://webpack.js.org/plugins/split-chunks-plugin/
 *
 */

module.exports = {
  entry: { bundle: './src/index.ts' },
	output: {
		filename: 'static/[name].js',
		chunkFilename: 'static/[name].[contenthash].js',
		path: path.resolve(__dirname, '..', 'html')
	},
	plugins: [new webpack.ProgressPlugin()],
	module: {
		rules: [
			{
				test: /.(ts|tsx)?$/,
				loader: 'ts-loader',
				include: [path.resolve(__dirname, 'src')],
				exclude: [/node_modules/]
			}
		]
	},
 	optimization: {
 		splitChunks: {
 			cacheGroups: {
 				vendors: {
 					priority: -10,
 					test: /[\\/]node_modules[\\/]/,
                                        chunks: 'all',
 				}
 			},
 
 			chunks: 'async',
 			minChunks: 1,
 			minSize: 30000,
 			name: true
 		}
 	},
	resolve: {
		extensions: ['.tsx', '.ts', '.js']
	}
};
