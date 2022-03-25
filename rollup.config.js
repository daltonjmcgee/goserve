import scss from 'rollup-plugin-scss'

export default {
  input: 'web/scripts/index.js',
  output: {
    file: 'web/public/static/main.js',
    format: 'cjs'
  },
  plugins: [
    scss({
      output: 'web/public/static/style.css',
      watch: ['web/styles']
    })
  ]
};