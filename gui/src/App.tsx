import { useState, useEffect } from "react";
import Layout from "./components/Layout";
import SkillList from "./components/SkillList";
import SkillDetail from "./components/SkillDetail";
import EngagementPicker from "./components/EngagementPicker";
import ExecutionPanel from "./components/ExecutionPanel";
import PhaseNav from "./components/PhaseNav";
import ArtifactBrowser from "./components/ArtifactBrowser";
import ChecklistView from "./components/ChecklistView";
import AssistantPanel from "./components/AssistantPanel";
import { useSkills } from "./hooks/useSkills";
import { useEngagements } from "./hooks/useEngagements";
import { useExecution } from "./hooks/useExecution";
import { usePhase } from "./hooks/usePhase";
import type { SkillDetail as SkillDetailType, AppView } from "./types";
import styles from "./App.module.css";

export default function App() {
  const { skills, loading: skillsLoading } = useSkills();
  const { engagements, loading: engagementsLoading } = useEngagements();
  const execution = useExecution();
  const { getSkillDetail } = useSkills();

  const [selectedSkills, setSelectedSkills] = useState<Set<string>>(new Set());
  const [selectedEngagement, setSelectedEngagement] = useState<string | null>(
    null
  );
  const [viewingSkill, setViewingSkill] = useState<string | null>(null);
  const [skillContent, setSkillContent] = useState<SkillDetailType | null>(
    null
  );
  const [view, setView] = useState<AppView>("detail");

  const { phase: phaseInfo } = usePhase(selectedEngagement);

  // Auto-set default view when phase changes
  useEffect(() => {
    if (!phaseInfo) return;
    switch (phaseInfo.phase) {
      case "pre-engagement":
        setView("detail");
        break;
      case "live":
        setView("detail");
        break;
      case "leave-behind":
        setView("artifacts");
        break;
    }
  }, [phaseInfo?.phase, selectedEngagement]);

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

      {selectedEngagement && phaseInfo && (
        <PhaseNav
          phase={phaseInfo.phase}
          currentView={view}
          onNavigate={setView}
          artifactCounts={phaseInfo.artifactCounts}
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

        {execution.status !== "idle" && view !== "execution" && (
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

  const main = (() => {
    switch (view) {
      case "execution":
        return (
          <ExecutionPanel
            messages={execution.messages}
            status={execution.status}
            onSendFollowUp={execution.sendFollowUp}
            onReset={() => {
              execution.reset();
              setView("detail");
            }}
          />
        );
      case "artifacts":
        return selectedEngagement ? (
          <ArtifactBrowser slug={selectedEngagement} />
        ) : null;
      case "checklists":
        return selectedEngagement ? (
          <ChecklistView slug={selectedEngagement} />
        ) : null;
      case "assistant":
        return selectedEngagement ? (
          <AssistantPanel slug={selectedEngagement} />
        ) : null;
      case "detail":
      default:
        return skillContent ? (
          <SkillDetail
            name={skillContent.name}
            content={skillContent.content}
          />
        ) : null;
    }
  })();

  return <Layout sidebar={sidebar} main={main} />;
}
