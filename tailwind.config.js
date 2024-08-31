/** @type {import('tailwindcss').Config} */
export default {
  content: ["./infra/http/htmx/**/*.{go,js,templ,html}"],
  theme: {
    fontFamily: {
      'sans': ['Outfit', 'sans-serif'],
    },
    extend: {},
  },
  plugins: [require('daisyui')],
  daisyui: {
    // DaisyUI theme builder: https://themes.ionevolve.com/
    themes: [
      {
        'light': {
          'primary': '#F06548',
          'primary-focus': '#ed3f1d',
          'primary-content': '#ffffff',

          'secondary': '#585f4c',
          'secondary-focus': '#3F4436',
          'secondary-content': '#ffffff',

          'accent': '#b58d69',
          'accent-focus': '#896443',
          'accent-content': '#ffffff',

          'neutral': '#3b424e',
          'neutral-focus': '#2a2e37',
          'neutral-content': '#ffffff',

          'base-100': '#ffffff',
          'base-200': '#f9fafb',
          'base-300': '#ced3d9',
          'base-content': '#1e2734',

          'info': '#1c92f2',
          'success': '#009485',
          'warning': '#ff9900',
          'error': '#ff5724',

          '--rounded-box': '1rem',
          '--rounded-btn': '.5rem',
          '--rounded-badge': '1.9rem',

          '--animation-btn': '.25s',
          '--animation-input': '.2s',

          '--btn-text-case': 'uppercase',
          '--navbar-padding': '.5rem',
          '--border-btn': '1px',
        },
      },
    ],
  }
}

