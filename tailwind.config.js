/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{html,js}",
    "./internal/**/*.templ",
    "./templates/**/*.go",
  ],
  theme: {
    fontFamily: {
      main: ["Plus Jakarta Sans", "sans-serif"],
    },
    extend: {},
  },
  plugins: [],
};
