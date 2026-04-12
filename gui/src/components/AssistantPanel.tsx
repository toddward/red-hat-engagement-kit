import { useRef, useEffect, useState } from "react";
import StreamOutput from "./StreamOutput";
import TypingIndicator from "./TypingIndicator";
import { useAssistant } from "../hooks/useAssistant";
import styles from "./AssistantPanel.module.css";

interface AssistantPanelProps {
  slug: string;
}

export default function AssistantPanel({ slug }: AssistantPanelProps) {
  const { messages, status, askQuestion, reset } = useAssistant(slug);
  const scrollRef = useRef<HTMLDivElement>(null);
  const [input, setInput] = useState("");

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = input.trim();
    if (!trimmed || status === "running") return;
    askQuestion(trimmed);
    setInput("");
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <span className={styles.title}>Engagement Assistant</span>
          {status === "running" && (
            <span className={`${styles.badge} ${styles.badgeRunning}`}>
              <span className={styles.spinner} />
              Thinking
            </span>
          )}
          {status === "done" && (
            <span className={`${styles.badge} ${styles.badgeDone}`}>Done</span>
          )}
        </div>
        {status === "done" && (
          <button onClick={reset} className={styles.clearButton}>
            Clear
          </button>
        )}
      </div>

      <div ref={scrollRef} className={styles.scrollArea}>
        {messages.length === 0 ? (
          <div className={styles.emptyHint}>
            Ask questions about the engagement artifacts.
            <div className={styles.emptySubtitle}>
              Answers are grounded in CONTEXT.md, discovery reports, assessments,
              and deliverables.
            </div>
          </div>
        ) : (
          <>
            <StreamOutput messages={messages} />
            {status === "running" && <TypingIndicator />}
          </>
        )}
      </div>

      <form onSubmit={handleSubmit} className={styles.inputBar}>
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="What would you like to know about this engagement?"
          disabled={status === "running"}
          className={styles.textInput}
        />
        <button
          type="submit"
          disabled={status === "running" || !input.trim()}
          className={styles.askButton}
        >
          Ask
        </button>
      </form>
    </div>
  );
}
