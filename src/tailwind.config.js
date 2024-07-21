/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['../views/**/*.{html,js}'],
    theme: {
        extend: {
            fontFamily: {
                "mono": ["Berkeley Mono", "monospace"],
            }
        },
    },
    plugins: [require('daisyui')],
    daisyui: {
        themes: ['cupcake', 'dracula'],
        styled: true,
        base: true,
        utils: true,
        darkTheme: "dracula",
    },
}
