/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        bg: {
          DEFAULT: "#0f1115",
          soft: "#12151b"
        },
        primary: {
          DEFAULT: "#7c3aed"
        },
        accent: {
          DEFAULT: "#22d3ee"
        }
      },
      boxShadow: {
        soft: "0 2px 12px rgba(0,0,0,0.35)"
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "Segoe UI", "Arial", "sans-serif"]
      }
    }
  },
  plugins: []
};
