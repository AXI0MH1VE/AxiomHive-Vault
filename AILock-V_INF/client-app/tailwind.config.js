// SUPRAPROTOCOL V∞: State-of-the-Art UX/UI - Tailwind CSS/Atomic Design Setup
/** @type {import('tailwindcss').Config} */
module.exports = {
  // Use custom theme colors for modern, professional look
  darkMode: ["class"],
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './src/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        // Custom Auth System Palette (inspired by zero-trust/dark theme)
        'auth-primary': 'hsl(210 40% 96.1%)', // Light
        'auth-secondary': 'hsl(214.3 31.8% 91.4%)', // Muted
        'auth-accent': 'hsl(222.2 47.4% 11.2%)', // Background
        'auth-foreground': 'hsl(210 40% 98%)', // Foreground text
        'danger-zone': 'hsl(0 84.2% 60.2%)', // Safety and audit alerts
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
}
