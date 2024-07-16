/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/*.html"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
}

