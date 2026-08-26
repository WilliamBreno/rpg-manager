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
        // Tokens semânticos do shadcn/ui, mapeados para as variáveis CSS em
        // index.css (que por sua vez já apontam pra paleta dourada acima).
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
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
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      keyframes: {
        'rpg-fade-in': {
          '0%':   { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'rpg-dice-entrance': {
          '0%':   { opacity: '0', transform: 'rotate(-140deg) scale(0.4)' },
          '60%':  { opacity: '1', transform: 'rotate(12deg) scale(1.08)' },
          '100%': { opacity: '1', transform: 'rotate(0deg) scale(1)' },
        },
        'rpg-rise': {
          '0%':   { opacity: '0', transform: 'translateY(14px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'rpg-level-pop': {
          '0%':   { opacity: '0', transform: 'scale(0.3)' },
          '55%':  { opacity: '1', transform: 'scale(1.15)' },
          '80%':  { transform: 'scale(0.96)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
        'rpg-ring-expand': {
          '0%':   { opacity: '0.9', transform: 'scale(0.4)' },
          '100%': { opacity: '0', transform: 'scale(2.2)' },
        },
        'rpg-burst-particle': {
          '0%':   { opacity: '1', transform: 'translate(0, 0) scale(1)' },
          '100%': { opacity: '0', transform: 'translate(var(--dx), var(--dy)) scale(0.4)' },
        },
      },
      animation: {
        'rpg-fade-in':       'rpg-fade-in 0.4s ease-out both',
        'rpg-dice-entrance': 'rpg-dice-entrance 0.9s cubic-bezier(0.34,1.56,0.64,1) both',
        'rpg-rise':          'rpg-rise 0.6s ease-out both',
        'rpg-level-pop':     'rpg-level-pop 0.7s cubic-bezier(0.34,1.56,0.64,1) both',
        'rpg-ring-expand':   'rpg-ring-expand 1.1s ease-out both',
        'rpg-burst-particle':'rpg-burst-particle 0.9s ease-out both',
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
}