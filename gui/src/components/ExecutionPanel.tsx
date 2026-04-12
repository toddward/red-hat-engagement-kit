import { useRef, useEffect, useState } from "react";
import StreamOutput from "./StreamOutput";
import TypingIndicator from "./TypingIndicator";
import type { StreamMessage, ExecutionStatus } from "../types";
import styles from "./ExecutionPanel.module.css";

interface ExecutionPanelProps {
  messages: StreamMessage[];
  status: ExecutionStatus;
  onSendFollowUp: (message: string) => void;
  onReset: () => void;
}

export default function ExecutionPanel({
  messages,
  status,
  onSendFollowUp,
  onReset,
}: ExecutionPanelProps) {
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
    if (!trimmed) return;
    onSendFollowUp(trimmed);
    setInput("");
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <span className={styles.title}>Execution Output</span>
          <StatusBadge status={status} />
        </div>
        {(status === "done" || status === "error") && (
          <button onClick={onReset} className={styles.clearButton}>
            Clear
          </button>
        )}
      </div>

      <div ref={scrollRef} className={styles.scrollArea}>
        {messages.length === 0 ? (
          <div className={styles.emptyHint}>
            Select skills and click Execute to begin.
          </div>
        ) : (
          <>
            <StreamOutput messages={messages} />
            {status === "running" && <TypingIndicator />}
          </>
        )}
      </div>

      {(status === "running" || status === "waiting" || status === "done") && (
        <form onSubmit={handleSubmit} className={styles.inputBar}>
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a response..."
            disabled={status === "running"}
            className={styles.textInput}
          />
          <button
            type="submit"
            disabled={status === "running" || !input.trim()}
            className={styles.sendButton}
          >
            Send
          </button>
        </form>
      )}
    </div>
  );
}

function StatusBadge({ status }: { status: ExecutionStatus }) {
  const badgeClass: Record<ExecutionStatus, string> = {
    idle: styles.badgeIdle,
    running: styles.badgeRunning,
    waiting: styles.badgeWaiting,
    done: styles.badgeDone,
    error: styles.badgeError,
  };

  const labels: Record<ExecutionStatus, string> = {
    idle: "Idle",
    running: "Running",
    waiting: "Waiting for input",
    done: "Done",
    error: "Error",
  };

  return (
    <span className={`${styles.badge} ${badgeClass[status]}`}>
      {status === "running" && <span className={styles.spinner} />}
      {labels[status]}
    </span>
  );
}
