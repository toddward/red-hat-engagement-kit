import type { Engagement } from "../types";
import styles from "./EngagementPicker.module.css";

interface EngagementPickerProps {
  engagements: Engagement[];
  selected: string | null;
  onSelect: (slug: string | null) => void;
}

export default function EngagementPicker({
  engagements,
  selected,
  onSelect,
}: EngagementPickerProps) {
  return (
    <div className={styles.container}>
      <h2 className={styles.sectionLabel}>Engagement</h2>
      <select
        value={selected ?? ""}
        onChange={(e) => onSelect(e.target.value || null)}
        className={styles.select}
      >
        <option value="">(New engagement)</option>
        {engagements.map((eng) => (
          <option key={eng.slug} value={eng.slug}>
            {eng.slug}
            {eng.hasContext ? "" : " (no context)"}
          </option>
        ))}
      </select>
    </div>
  );
}
