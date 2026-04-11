import { useState, useEffect } from "react";
import Layout from "./components/Layout";
import SkillList from "./components/SkillList";
import SkillDetail from "./components/SkillDetail";
import EngagementPicker from "./components/EngagementPicker";
import ExecutionPanel from "./components/ExecutionPanel";
import { useSkills } from "./hooks/useSkills";
import { useEngagements } from "./hooks/useEngagements";
import { useExecution } from "./hooks/useExecution";
import type { SkillDetail as SkillDetailType } from "./types";
import styles from "./App.module.css";

export default function App() {
  const { skills, loading: skillsLoading } = useSkills();
  const { engagements, loading: engagementsLoading } = useEngagements();
  const execution = useExecution();

  const [selectedSkills, setSelectedSkills] = useState<Set<string>>(new Set());
  const [selectedEngagement, setSelectedEngagement] = useState<string | null>(
    null
  );
  const [viewingSkill, setViewingSkill] = useState<string | null>(null);
  const [skillContent, setSkillContent] = useState<SkillDetailType | null>(
    null
  );
  const [view, setView] = useState<"detail" | "execution">("detail");
  const { getSkillDetail } = useSkills();

  useEffect(() => {
    if (!viewingSkill) {
      setSkillContent(null);
      return;
    }
    getSkillDetail(viewingSkill).then((detail) => {
      if (detail) setSkillContent(detail);
    });
  }, [viewingSkill, getSkillDetail]);

  const handleToggleSkill = (name: string) => {
    setSelectedSkills((prev) => {
      const next = new Set(prev);
      if (next.has(name)) {
        next.delete(name);
      } else {
        next.add(name);
      }
      return next;
    });
  };

  const isDisabled =
    selectedSkills.size === 0 || execution.status === "running";

  const handleExecute = () => {
    if (isDisabled) return;
    setView("execution");
    execution.startExecution(
      Array.from(selectedSkills),
      selectedEngagement ?? undefined
    );
  };

  const sidebar = (
    <>
      {skillsLoading ? (
        <div className={styles.loadingText}>Loading skills...</div>
      ) : (
        <SkillList
          skills={skills}
          selected={selectedSkills}
          onToggle={handleToggleSkill}
          onView={(name) => {
            setViewingSkill(name);
            setView("detail");
          }}
          activeView={view === "detail" ? viewingSkill : null}
        />
      )}

      {engagementsLoading ? (
        <div className={styles.loadingText} style={{ marginTop: 24 }}>
          Loading engagements...
        </div>
      ) : (
        <EngagementPicker
          engagements={engagements}
          selected={selectedEngagement}
          onSelect={setSelectedEngagement}
        />
      )}

      <div className={styles.actionArea}>
        <button
          onClick={handleExecute}
          disabled={isDisabled}
          className={`${styles.executeButton} ${isDisabled ? styles.executeButtonDisabled : ""}`}
        >
          {execution.status === "running"
            ? "Running..."
            : `Execute ${selectedSkills.size > 0 ? `(${selectedSkills.size})` : ""}`}
        </button>

        {execution.status !== "idle" && view === "detail" && (
          <button
            onClick={() => setView("execution")}
            className={styles.viewOutputButton}
          >
            View execution output
          </button>
        )}
      </div>
    </>
  );

  const main =
    view === "execution" ? (
      <ExecutionPanel
        messages={execution.messages}
        status={execution.status}
        onSendFollowUp={execution.sendFollowUp}
        onReset={() => {
          execution.reset();
          setView("detail");
        }}
      />
    ) : skillContent ? (
      <SkillDetail name={skillContent.name} content={skillContent.content} />
    ) : null;

  return <Layout sidebar={sidebar} main={main} />;
}
