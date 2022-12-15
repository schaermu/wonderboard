/** @type {import('tailwindcss').Config} */

const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{svelte,ts}"
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          'Fira Sans',
          ...defaultTheme.fontFamily.sans,
        ],
        serif: [
          'Bree Serif',
          defaultTheme.fontFamily.serif,
        ]
      },
    },
  },
  plugins: [],
}
