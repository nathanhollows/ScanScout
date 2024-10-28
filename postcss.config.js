// postcss.config.js
module.exports = {
  plugins: [
    require('postcss-import'), // Enables `@import` for CSS files
    require('tailwindcss'),
    require('autoprefixer'),
  ]
}
