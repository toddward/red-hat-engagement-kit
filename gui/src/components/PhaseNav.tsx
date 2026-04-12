import type { EngagementPhase, PhaseInfo, AppView } from "../types";
import styles from "./PhaseNav.module.css";

interface PhaseNavProps {
  phase: EngagementPhase;
  currentView: AppView;
  onNavigate: (view: AppView) => void;
  artifactCounts: PhaseInfo["artifactCounts"];
}

interface NavEntry {
  view: AppView;
  label: string;
  icon: string;
  count?: number;
  prominence: "primary" | "secondary" | "dimmed";
}

function getNavEntries(
  phase: EngagementPhase,
  counts: PhaseInfo["artifactCounts"]
): NavEntry[] {
  const totalArtifacts =
    counts.discovery + counts.assessments + counts.deliverables;

  switch (phase) {
    case "pre-engagement":
      return [
        { view: "checklists", label: "Checklists", icon: "\u2611", count: counts.checklists, prominence: "primary" },
        { view: "detail", label: "Skills", icon: "\u2699", prominence: "secondary" },
        { view: "artifacts", label: "Artifacts", icon: "\u2630", count: totalArtifacts, prominence: "dimmed" },
      ];
    case "live":
      return [
        { view: "detail", label: "Skills", icon: "\u2699", prominence: "primary" },
        { view: "artifacts", label: "Artifacts", icon: "\u2630", count: totalArtifacts, prominence: "secondary" },
        { view: "checklists", label: "Checklists", icon: "\u2611", count: counts.checklists, prominence: "secondary" },
        { view: "assistant", label: "Assistant", icon: "\u2709", prominence: "secondary" },
      ];
    case "leave-behind":
      return [
        { view: "artifacts", label: "Artifacts", icon: "\u2630", count: totalArtifacts, prominence: "primary" },
        { view: "assistant", label: "Assistant", icon: "\u2709", prominence: "primary" },
        { view: "detail", label: "Skills", icon: "\u2699", prominence: "secondary" },
        { view: "checklists", label: "Checklists", icon: "\u2611", count: counts.checklists, prominence: "dimmed" },
      ];
  }
}

const phaseLabels: Record<EngagementPhase, string> = {
  "pre-engagement": "Pre-engagement",
  live: "Live Engagement",
  "leave-behind": "Leave-behind",
};

const phaseStyles: Record<EngagementPhase, string> = {
  "pre-engagement": styles.phasePre,
  live: styles.phaseLive,
  "leave-behind": styles.phaseLeave,
};

export default function PhaseNav({
  phase,
  currentView,
  onNavigate,
  artifactCounts,
}: PhaseNavProps) {
  const entries = getNavEntries(phase, artifactCounts);

  return (
    <div className={styles.container}>
      <div className={`${styles.phaseBadge} ${phaseStyles[phase]}`}>
        <span className={styles.phaseDot} />
        {phaseLabels[phase]}
      </div>

      <div className={styles.navList}>
        {entries.map((entry) => {
          const isActive = currentView === entry.view;
          const isDimmed = entry.prominence === "dimmed";

          const itemClass = [
            styles.navItem,
            isActive ? styles.navItemActive : "",
            isDimmed ? styles.navItemDimmed : "",
          ]
            .filter(Boolean)
            .join(" ");

          return (
            <button
              key={entry.view}
              className={itemClass}
              onClick={() => !isDimmed && onNavigate(entry.view)}
              disabled={isDimmed}
            >
              <span className={styles.navIcon}>{entry.icon}</span>
              <span className={styles.navLabel}>{entry.label}</span>
              {entry.count != null && entry.count > 0 && (
                <span className={styles.navCount}>{entry.count}</span>
              )}
            </button>
          );
        })}
      </div>
    </div>
  );
}
