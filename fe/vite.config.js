import { createHtmlPlugin } from "vite-plugin-html";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [
    createHtmlPlugin({
      minify: true, // Enable minification
      minifyOptions: {
        collapseWhitespace: true, // Remove whitespaces
        removeComments: true, // Remove comments
        removeRedundantAttributes: true, // Remove redundant attributes
        useShortDoctype: true, // Use short doctype
        removeEmptyAttributes: true, // Remove empty attributes
        removeOptionalTags: true, // Remove optional tags
      },
    }),
  ],
});
