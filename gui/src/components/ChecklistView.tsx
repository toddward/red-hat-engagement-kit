import { useState } from "react";
import { useChecklists } from "../hooks/useChecklists";
import type { Checklist } from "../types";
import styles from "./ChecklistView.module.css";

interface ChecklistViewProps {
  slug: string;
}

export default function ChecklistView({ slug }: ChecklistViewProps) {
  const { checklists, loading, instantiate, toggleItem, summary } =
    useChecklists(slug);

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.subtitle}>Loading checklists...</div>
      </div>
    );
  }

  if (checklists.length === 0) {
    return (
      <div className={styles.container}>
        <div className={styles.emptyState}>
          <div className={styles.emptyTitle}>No Checklists Yet</div>
          <div className={styles.emptyDesc}>
            Initialize role-based preparation checklists from templates.
          </div>
          <button onClick={instantiate} className={styles.instantiateButton}>
            Initialize Checklists
          </button>
        </div>
      </div>
    );
  }

  const progressColor =
    summary.percent >= 80
      ? "var(--rh-green)"
      : summary.percent >= 40
        ? "var(--rh-yellow)"
        : "var(--rh-red)";

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.title}>Engagement Readiness</div>
        <div className={styles.subtitle}>
          {summary.checked} of {summary.total} items completed
        </div>
      </div>

      <div className={styles.summary}>
        <div className={styles.summaryLabel}>Overall</div>
        <div className={styles.progressTrack}>
          <div
            className={styles.progressFill}
            style={{
              width: `${summary.percent}%`,
              background: progressColor,
            }}
          />
        </div>
        <div className={styles.summaryPercent} style={{ color: progressColor }}>
          {summary.percent}%
        </div>
      </div>

      {checklists.map((cl) => (
        <ChecklistCard
          key={cl.fileName}
          checklist={cl}
          onToggle={(line, checked) => toggleItem(cl.fileName, line, checked)}
        />
      ))}
    </div>
  );
}

function ChecklistCard({
  checklist,
  onToggle,
}: {
  checklist: Checklist;
  onToggle: (line: number, checked: boolean) => void;
}) {
  const [expanded, setExpanded] = useState(true);

  const barColor =
    checklist.completionPercent >= 80
      ? "var(--rh-green)"
      : checklist.completionPercent >= 40
        ? "var(--rh-yellow)"
        : "var(--rh-red)";

  return (
    <div className={styles.checklistCard}>
      <div
        className={styles.checklistHeader}
        onClick={() => setExpanded((p) => !p)}
      >
        <div className={styles.checklistName}>
          <span
            className={`${styles.checklistChevron} ${expanded ? styles.checklistChevronOpen : ""}`}
          >
            &#9656;
          </span>
          {checklist.name}
        </div>
        <div className={styles.checklistMiniBar}>
          <div
            className={styles.checklistMiniBarFill}
            style={{
              width: `${checklist.completionPercent}%`,
              background: barColor,
            }}
          />
        </div>
      </div>

      {expanded && (
        <div className={styles.checklistBody}>
          {checklist.sections.map((section) => (
            <div key={section.title}>
              <div className={styles.sectionTitle}>{section.title}</div>
              {section.items.map((item) => (
                <div
                  key={item.line}
                  className={styles.item}
                  onClick={() => onToggle(item.line, !item.checked)}
                >
                  <input
                    type="checkbox"
                    checked={item.checked}
                    onChange={() => onToggle(item.line, !item.checked)}
                    onClick={(e) => e.stopPropagation()}
                    className={styles.itemCheckbox}
                  />
                  <span
                    className={`${styles.itemText} ${item.checked ? styles.itemTextChecked : ""}`}
                  >
                    {item.text}
                  </span>
                </div>
              ))}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
