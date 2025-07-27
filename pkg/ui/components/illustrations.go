package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// LoginIllustration renders a modern login illustration with characters
func LoginIllustration() Node {
	return Div(
		Class("flex justify-center mb-8"),
		Div(
			Class("w-72 h-48 sm:w-80 sm:h-52"),
			// Modern character-based login illustration
			Raw(`<svg viewBox="0 0 400 300" fill="none" xmlns="http://www.w3.org/2000/svg">
				<!-- Background elements -->
				<circle cx="80" cy="50" r="25" fill="#E0F2FE" opacity="0.6"/>
				<circle cx="320" cy="80" r="20" fill="#F3E8FF" opacity="0.5"/>
				<rect x="50" y="220" width="30" height="30" rx="15" fill="#FEF3C7" opacity="0.4"/>
				
				<!-- Main character - person with laptop -->
				<g transform="translate(180, 80)">
					<!-- Body -->
					<ellipse cx="0" cy="60" rx="35" ry="45" fill="#3B82F6"/>
					<!-- Head -->
					<circle cx="0" cy="0" r="25" fill="#FBBF24"/>
					<!-- Hair -->
					<path d="M-20 -15 Q0 -30 20 -15 Q15 -25 0 -25 Q-15 -25 -20 -15" fill="#1F2937"/>
					<!-- Eyes -->
					<circle cx="-8" cy="-3" r="2" fill="#1F2937"/>
					<circle cx="8" cy="-3" r="2" fill="#1F2937"/>
					<!-- Smile -->
					<path d="M-8 8 Q0 15 8 8" stroke="#1F2937" stroke-width="2" fill="none"/>
					<!-- Arms -->
					<ellipse cx="-25" cy="40" rx="8" ry="20" fill="#FBBF24" transform="rotate(-20)"/>
					<ellipse cx="25" cy="40" rx="8" ry="20" fill="#FBBF24" transform="rotate(20)"/>
					<!-- Laptop -->
					<rect x="-20" y="35" width="40" height="25" rx="3" fill="#1F2937"/>
					<rect x="-18" y="37" width="36" height="15" fill="#10B981"/>
					<circle cx="0" cy="44" r="2" fill="#FFFFFF"/>
				</g>
				
				<!-- Security icons floating around -->
				<g transform="translate(100, 120)">
					<circle cx="0" cy="0" r="15" fill="#EF4444" opacity="0.8"/>
					<rect x="-6" y="-8" width="12" height="10" rx="2" fill="white"/>
					<path d="M-3 -8 Q0 -12 3 -8" stroke="white" stroke-width="2" fill="none"/>
				</g>
				
				<g transform="translate(300, 150)">
					<circle cx="0" cy="0" r="12" fill="#10B981" opacity="0.8"/>
					<path d="M-4 0 L-1 3 L4 -3" stroke="white" stroke-width="2" fill="none"/>
				</g>
				
				<!-- Decorative elements -->
				<path d="M50 280 Q200 250 350 280" stroke="#8B5CF6" stroke-width="3" fill="none" opacity="0.3"/>
				<circle cx="120" cy="260" r="4" fill="#F59E0B" opacity="0.6"/>
				<circle cx="280" cy="40" r="3" fill="#EF4444" opacity="0.5"/>
			</svg>`),
		),
	)
}

// RegisterIllustration renders a modern registration illustration with characters
func RegisterIllustration() Node {
	return Div(
		Class("flex justify-center mb-8"),
		Div(
			Class("w-72 h-48 sm:w-80 sm:h-52"),
			// Modern character-based register illustration
			Raw(`<svg viewBox="0 0 400 300" fill="none" xmlns="http://www.w3.org/2000/svg">
				<!-- Background elements -->
				<circle cx="70" cy="60" r="30" fill="#F0F9FF" opacity="0.7"/>
				<circle cx="330" cy="90" r="25" fill="#FDF4FF" opacity="0.6"/>
				<rect x="40" y="240" width="25" height="25" rx="12" fill="#ECFDF5" opacity="0.5"/>
				
				<!-- Three diverse characters -->
				<!-- Character 1 - Left -->
				<g transform="translate(80, 90)">
					<!-- Body -->
					<ellipse cx="0" cy="50" rx="30" ry="40" fill="#F59E0B"/>
					<!-- Head -->
					<circle cx="0" cy="0" r="22" fill="#FBBF24"/>
					<!-- Hair -->
					<path d="M-18 -12 Q0 -25 18 -12 Q12 -22 0 -22 Q-12 -22 -18 -12" fill="#7C2D12"/>
					<!-- Eyes -->
					<circle cx="-6" cy="-2" r="1.5" fill="#1F2937"/>
					<circle cx="6" cy="-2" r="1.5" fill="#1F2937"/>
					<!-- Smile -->
					<path d="M-6 6 Q0 12 6 6" stroke="#1F2937" stroke-width="1.5" fill="none"/>
					<!-- Arms raised -->
					<ellipse cx="-20" cy="30" rx="6" ry="18" fill="#FBBF24" transform="rotate(-45)"/>
					<ellipse cx="20" cy="30" rx="6" ry="18" fill="#FBBF24" transform="rotate(45)"/>
					<!-- Legs -->
					<ellipse cx="-10" cy="80" rx="6" ry="15" fill="#DC2626"/>
					<ellipse cx="10" cy="80" rx="6" ry="15" fill="#DC2626"/>
				</g>
				
				<!-- Character 2 - Center with laptop -->
				<g transform="translate(200, 100)">
					<!-- Body -->
					<ellipse cx="0" cy="45" rx="32" ry="42" fill="#8B5CF6"/>
					<!-- Head -->
					<circle cx="0" cy="0" r="24" fill="#FDE68A"/>
					<!-- Hair -->
					<path d="M-20 -10 Q0 -28 20 -10 Q15 -25 0 -25 Q-15 -25 -20 -10" fill="#059669"/>
					<!-- Glasses -->
					<circle cx="-8" cy="-2" r="6" fill="none" stroke="#1F2937" stroke-width="1.5"/>
					<circle cx="8" cy="-2" r="6" fill="none" stroke="#1F2937" stroke-width="1.5"/>
					<path d="M2 -2 L6 -2" stroke="#1F2937" stroke-width="1.5"/>
					<!-- Eyes -->
					<circle cx="-8" cy="-2" r="1.5" fill="#1F2937"/>
					<circle cx="8" cy="-2" r="1.5" fill="#1F2937"/>
					<!-- Smile -->
					<path d="M-8 8 Q0 15 8 8" stroke="#1F2937" stroke-width="1.5" fill="none"/>
					<!-- Arms typing -->
					<ellipse cx="-22" cy="35" rx="6" ry="16" fill="#FDE68A" transform="rotate(-10)"/>
					<ellipse cx="22" cy="35" rx="6" ry="16" fill="#FDE68A" transform="rotate(10)"/>
					<!-- Laptop -->
					<rect x="-18" y="30" width="36" height="20" rx="2" fill="#1F2937"/>
					<rect x="-16" y="32" width="32" height="12" fill="#10B981"/>
				</g>
				
				<!-- Character 3 - Right -->
				<g transform="translate(320, 95)">
					<!-- Body -->
					<ellipse cx="0" cy="48" rx="28" ry="38" fill="#EC4899"/>
					<!-- Head -->
					<circle cx="0" cy="0" r="20" fill="#FED7AA"/>
					<!-- Hair -->
					<circle cx="0" cy="-5" r="22" fill="#1F2937"/>
					<!-- Eyes -->
					<circle cx="-5" cy="-1" r="1.5" fill="#1F2937"/>
					<circle cx="5" cy="-1" r="1.5" fill="#1F2937"/>
					<!-- Smile -->
					<path d="M-5 7 Q0 12 5 7" stroke="#1F2937" stroke-width="1.5" fill="none"/>
					<!-- Arms waving -->
					<ellipse cx="-18" cy="28" rx="5" ry="15" fill="#FED7AA" transform="rotate(-30)"/>
					<ellipse cx="18" cy="28" rx="5" ry="15" fill="#FED7AA" transform="rotate(30)"/>
					<!-- Legs -->
					<ellipse cx="-8" cy="75" rx="5" ry="12" fill="#1E40AF"/>
					<ellipse cx="8" cy="75" rx="5" ry="12" fill="#1E40AF"/>
				</g>
				
				<!-- Plus icon for joining -->
				<g transform="translate(140, 60)">
					<circle cx="0" cy="0" r="12" fill="#10B981"/>
					<path d="M-6 0 L6 0 M0 -6 L0 6" stroke="white" stroke-width="2"/>
				</g>
				
				<!-- Welcome elements -->
				<g transform="translate(260, 50)">
					<circle cx="0" cy="0" r="8" fill="#F59E0B"/>
					<path d="M-3 -1 L0 2 L3 -2" stroke="white" stroke-width="1.5" fill="none"/>
				</g>
				
				<!-- Decorative bottom wave -->
				<path d="M0 270 Q100 250 200 260 Q300 270 400 255" stroke="#3B82F6" stroke-width="3" fill="none" opacity="0.4"/>
				<circle cx="150" cy="280" r="3" fill="#8B5CF6" opacity="0.6"/>
				<circle cx="250" cy="275" r="4" fill="#EC4899" opacity="0.5"/>
			</svg>`),
		),
	)
}

// ForgotPasswordIllustration renders a modern forgot password illustration
func ForgotPasswordIllustration() Node {
	return Div(
		Class("flex justify-center mb-6"),
		Div(
			Class("w-32 h-32 sm:w-40 sm:h-40"),
			// Modern forgot password illustration SVG
			Raw(`<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
				<!-- Background circle with gradient -->
				<defs>
					<linearGradient id="forgotGrad" x1="0%" y1="0%" x2="100%" y2="100%">
						<stop offset="0%" style="stop-color:#F59E0B;stop-opacity:0.2" />
						<stop offset="100%" style="stop-color:#EF4444;stop-opacity:0.2" />
					</linearGradient>
				</defs>
				<circle cx="100" cy="100" r="90" fill="url(#forgotGrad)" />
				
				<!-- Email envelope -->
				<rect x="70" y="85" width="60" height="40" rx="6" fill="white" stroke="#F59E0B" stroke-width="2"/>
				<path d="M70 85l30 20 30-20" stroke="#F59E0B" stroke-width="2" fill="none"/>
				
				<!-- Key icon -->
				<circle cx="85" cy="70" r="8" fill="#EF4444"/>
				<rect x="93" y="66" width="15" height="8" rx="2" fill="#EF4444"/>
				<rect x="105" y="68" width="3" height="4" fill="white"/>
				<rect x="105" y="69" width="5" height="2" fill="white"/>
				
				<!-- Question mark -->
				<circle cx="120" cy="130" r="10" fill="#64748B"/>
				<path d="M116 127c0-2 2-3 4-3s4 1 4 3-2 2-4 3v2" stroke="white" stroke-width="1.5" fill="none"/>
				<circle cx="120" cy="135" r="1" fill="white"/>
				
				<!-- Floating elements -->
				<circle cx="140" cy="70" r="3" fill="#F59E0B" opacity="0.6"/>
				<circle cx="60" cy="110" r="2" fill="#EF4444" opacity="0.4"/>
			</svg>`),
		),
	)
}

// EmailVerificationIllustration renders a modern email verification illustration
func EmailVerificationIllustration() Node {
	return Div(
		Class("flex justify-center mb-6"),
		Div(
			Class("w-40 h-40 sm:w-48 sm:h-48"),
			// Modern email verification illustration SVG
			Raw(`<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
				<!-- Background circle with gradient -->
				<defs>
					<linearGradient id="verifyGrad" x1="0%" y1="0%" x2="100%" y2="100%">
						<stop offset="0%" style="stop-color:#10B981;stop-opacity:0.2" />
						<stop offset="100%" style="stop-color:#3B82F6;stop-opacity:0.2" />
					</linearGradient>
				</defs>
				<circle cx="100" cy="100" r="90" fill="url(#verifyGrad)" />
				
				<!-- Email envelope -->
				<rect x="65" y="80" width="70" height="45" rx="8" fill="white" stroke="#10B981" stroke-width="3"/>
				<path d="M65 80l35 25 35-25" stroke="#10B981" stroke-width="3" fill="none"/>
				
				<!-- Checkmark -->
				<circle cx="100" cy="140" r="15" fill="#10B981"/>
				<path d="M92 140l5 5 10-10" stroke="white" stroke-width="3" fill="none"/>
				
				<!-- Sparkles -->
				<g fill="#F59E0B" opacity="0.7">
					<path d="M140 90l2-6 2 6 6-2-6 2z"/>
					<path d="M60 100l1.5-4 1.5 4 4-1.5-4 1.5z"/>
					<path d="M130 140l1-3 1 3 3-1-3 1z"/>
				</g>
				
				<!-- Floating elements -->
				<circle cx="50" cy="70" r="2" fill="#10B981" opacity="0.5"/>
				<circle cx="150" cy="120" r="3" fill="#3B82F6" opacity="0.4"/>
			</svg>`),
		),
	)
}

// ErrorIllustration renders a modern error illustration
func ErrorIllustration() Node {
	return Div(
		Class("flex justify-center mb-8"),
		Div(
			Class("w-48 h-48 sm:w-56 sm:h-56"),
			// Modern error illustration SVG
			Raw(`<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
				<!-- Background circle with gradient -->
				<defs>
					<linearGradient id="errorGrad" x1="0%" y1="0%" x2="100%" y2="100%">
						<stop offset="0%" style="stop-color:#EF4444;stop-opacity:0.15" />
						<stop offset="100%" style="stop-color:#F97316;stop-opacity:0.15" />
					</linearGradient>
				</defs>
				<circle cx="100" cy="100" r="90" fill="url(#errorGrad)" />
				
				<!-- Robot/character -->
				<rect x="80" y="70" width="40" height="35" rx="8" fill="#64748B"/>
				<circle cx="90" cy="85" r="4" fill="#EF4444"/>
				<circle cx="110" cy="85" r="4" fill="#EF4444"/>
				<rect x="95" y="92" width="10" height="3" rx="1.5" fill="#EF4444"/>
				
				<!-- Error symbol -->
				<circle cx="100" cy="125" r="12" fill="#EF4444"/>
				<path d="M95 120l10 10M105 120l-10 10" stroke="white" stroke-width="2"/>
				
				<!-- Broken elements -->
				<path d="M70 140l-5-5 5-5" stroke="#94A3B8" stroke-width="2" fill="none"/>
				<path d="M130 140l5-5-5-5" stroke="#94A3B8" stroke-width="2" fill="none"/>
				<rect x="75" y="148" width="8" height="8" rx="2" fill="#94A3B8" opacity="0.5"/>
				<rect x="117" y="145" width="6" height="6" rx="1" fill="#94A3B8" opacity="0.5"/>
			</svg>`),
		),
	)
}

// EmptyStateIllustration renders a modern empty state illustration
func EmptyStateIllustration() Node {
	return Div(
		Class("flex justify-center mb-8"),
		Div(
			Class("w-48 h-48 sm:w-56 sm:h-56"),
			// Modern empty state illustration SVG
			Raw(`<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
				<!-- Background circle with gradient -->
				<defs>
					<linearGradient id="emptyGrad" x1="0%" y1="0%" x2="100%" y2="100%">
						<stop offset="0%" style="stop-color:#8B5CF6;stop-opacity:0.1" />
						<stop offset="100%" style="stop-color:#3B82F6;stop-opacity:0.1" />
					</linearGradient>
				</defs>
				<circle cx="100" cy="100" r="90" fill="url(#emptyGrad)" />
				
				<!-- Empty box -->
				<rect x="70" y="90" width="60" height="45" rx="6" fill="none" stroke="#94A3B8" stroke-width="2" stroke-dasharray="5,5"/>
				
				<!-- Floating document icons -->
				<rect x="85" y="70" width="12" height="16" rx="2" fill="#E2E8F0" stroke="#94A3B8"/>
				<rect x="87" y="74" width="8" height="2" fill="#94A3B8"/>
				<rect x="87" y="78" width="6" height="2" fill="#94A3B8"/>
				
				<rect x="103" y="65" width="12" height="16" rx="2" fill="#E2E8F0" stroke="#94A3B8"/>
				<rect x="105" y="69" width="8" height="2" fill="#94A3B8"/>
				<rect x="105" y="73" width="6" height="2" fill="#94A3B8"/>
				
				<!-- Search icon -->
				<circle cx="100" cy="112" r="8" fill="none" stroke="#94A3B8" stroke-width="2"/>
				<path d="M106 118l4 4" stroke="#94A3B8" stroke-width="2"/>
				
				<!-- Floating elements -->
				<circle cx="130" cy="80" r="2" fill="#8B5CF6" opacity="0.3"/>
				<circle cx="70" cy="140" r="3" fill="#3B82F6" opacity="0.3"/>
				<rect x="140" y="130" width="4" height="4" rx="1" fill="#F59E0B" opacity="0.3"/>
			</svg>`),
		),
	)
}
