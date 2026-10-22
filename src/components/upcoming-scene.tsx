/**
 * Placeholder illustration for products that are announced but not on sale.
 * It deliberately reuses the tier-card language of the ChatGPT deck — same
 * 174×110 sheet, pastel gradients, border and shadow — so an upcoming product
 * sits in the grid without looking like a different design system.
 */
export function UpcomingScene({ name, index = 0 }: { name: string; index?: number }) {
  const letter = name.trim().slice(0, 1).toUpperCase()
  return <div className="upcoming-deck" data-tint={index % 4} aria-hidden="true">
    <span className="upcoming-sheet" data-depth="2" />
    <span className="upcoming-sheet" data-depth="1" />
    <span className="upcoming-sheet" data-depth="0">
      <span className="upcoming-sheet-label">即将上架</span>
      <strong className="upcoming-sheet-name">{name}</strong>
      <span className="upcoming-sheet-letter">{letter}</span>
    </span>
  </div>
}
