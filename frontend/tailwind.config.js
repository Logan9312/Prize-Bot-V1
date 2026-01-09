/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['Geist', 'system-ui', 'sans-serif']
			},
			colors: {
				surface: {
					900: '#111214',
					800: '#1a1b1e',
					700: '#232428',
					600: '#2e3035',
					500: '#3f4147',
					400: '#4e5058',
					300: '#6d6f78'
				},
				text: {
					primary: '#f2f3f5',
					secondary: '#b5bac1',
					muted: '#6d6f78'
				},
				accent: {
					DEFAULT: '#5865f2',
					hover: '#4752c4',
					light: '#7289da'
				},
				status: {
					success: '#23a559',
					warning: '#f0b232',
					danger: '#da373c'
				}
			},
			fontSize: {
				'fluid-xs': 'clamp(0.75rem, 0.7rem + 0.25vw, 0.875rem)',
				'fluid-sm': 'clamp(0.875rem, 0.8rem + 0.375vw, 1rem)',
				'fluid-base': 'clamp(1rem, 0.95rem + 0.25vw, 1.125rem)',
				'fluid-lg': 'clamp(1.125rem, 1rem + 0.5vw, 1.25rem)',
				'fluid-xl': 'clamp(1.25rem, 1.1rem + 0.75vw, 1.5rem)'
			},
			minHeight: {
				'touch': '44px'
			},
			minWidth: {
				'touch': '44px'
			},
			borderRadius: {
				DEFAULT: '8px'
			},
			animation: {
				'fade-in': 'fade-in 0.2s ease-out'
			},
			keyframes: {
				'fade-in': {
					'0%': { opacity: '0' },
					'100%': { opacity: '1' }
				}
			}
		}
	},
	plugins: []
};
