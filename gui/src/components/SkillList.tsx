import type { SkillSummary } from "../types";
import styles from "./SkillList.module.css";

interface SkillListProps {
  skills: SkillSummary[];
  selected: Set<string>;
  onToggle: (name: string) => void;
  onView: (name: string) => void;
  activeView: string | null;
}

export default function SkillList({
  skills,
  selected,
  onToggle,
  onView,
  activeView,
}: SkillListProps) {
  return (
    <div>
      <h2 className={styles.sectionLabel}>Skills</h2>
      <div className={styles.list}>
        {skills.map((skill) => {
          const isActive = activeView === skill.name;
          const isSelected = selected.has(skill.name);
          const cardClass = [
            styles.card,
            isActive ? styles.cardActive : "",
            isSelected ? styles.cardSelected : "",
          ]
            .filter(Boolean)
            .join(" ");

          return (
            <div
              key={skill.name}
              className={cardClass}
              onClick={() => onView(skill.name)}
            >
              <input
                type="checkbox"
                checked={isSelected}
                onChange={(e) => {
                  e.stopPropagation();
                  onToggle(skill.name);
                }}
                onClick={(e) => e.stopPropagation()}
                className={styles.checkbox}
              />
              <div className={styles.cardBody}>
                <div className={styles.skillName}>/{skill.name}</div>
                <div className={styles.skillDescription}>
                  {skill.description}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
