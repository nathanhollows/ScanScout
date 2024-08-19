/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/**/*.html", "./internal/templates/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/container-queries'),
    require("@tailwindcss/typography"),
    require("daisyui"),
  ],
  daisyui: {
    themes: ["cupcake", "dracula"],
    darkTheme: "dracula",
  },
};
