/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      // Sobrescreve gray para zinc (neutro, sem tom azulado do Tailwind padrão)
      // Assim bg-gray-900, bg-gray-800 etc. ficam com cara de RPG sem mudar cada arquivo
      colors: {
        gray: {
          50:  '#fafafa',
          100: '#f4f4f5',
          200: '#e4e4e7',
          300: '#d4d4d8',
          400: '#a1a1aa',
          500: '#71717a',
          600: '#52525b',
          700: '#3f3f46',
          800: '#27272a',
          900: '#18181b',
          950: '#09090b',
        },
        // Paleta dourada da logo
        rpg: {
          gold:          '#c9a84c',
          'gold-light':  '#e8c46a',
          'gold-dim':    '#8b6914',
          'gold-muted':  'rgba(201,168,76,0.15)',
          dark:          '#050505',
          surface:       '#111111',
          'surface-2':   '#1a1a1a',
          border:        '#2a2a2a',
          'border-gold': 'rgba(201,168,76,0.28)',
        },
      },
      fontFamily: {
        rpg: ['Georgia', '"Times New Roman"', 'serif'],
      },
      boxShadow: {
        'rpg-sm': '0 0 12px rgba(201,168,76,0.08)',
        'rpg':    '0 0 24px rgba(201,168,76,0.12)',
        'rpg-lg': '0 0 48px rgba(201,168,76,0.18)',
      },
      borderColor: {
        'rpg-gold': 'rgba(201,168,76,0.28)',
      },
    },
  },
  plugins: [],
}