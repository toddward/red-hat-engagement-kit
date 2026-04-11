import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { StreamMessage } from "../types";
import styles from "./StreamOutput.module.css";
import detailStyles from "./SkillDetail.module.css";

interface StreamOutputProps {
  messages: StreamMessage[];
}

export default function StreamOutput({ messages }: StreamOutputProps) {
  if (messages.length === 0) return null;

  return (
    <div className={styles.list}>
      {messages.map((msg) => (
        <MessageBlock key={msg.id} message={msg} />
      ))}
    </div>
  );
}

function MessageBlock({ message }: { message: StreamMessage }) {
  switch (message.type) {
    case "assistant":
      return (
        <div className={styles.assistantBlock}>
          <div className={detailStyles.markdownBody}>
            <Markdown remarkPlugins={[remarkGfm]}>
              {message.text ?? ""}
            </Markdown>
          </div>
        </div>
      );

    case "tool_use":
      return (
        <details className={styles.toolUseBlock}>
          <summary className={styles.toolUseSummary}>
            Tool: {message.tool}
          </summary>
          <pre className={styles.toolCode}>
            {typeof message.input === "string"
              ? message.input
              : JSON.stringify(message.input, null, 2)}
          </pre>
        </details>
      );

    case "tool_result":
      return (
        <details className={styles.toolResultBlock}>
          <summary className={styles.toolResultSummary}>
            Result: {message.tool}
          </summary>
          <pre className={styles.toolCode}>{message.output}</pre>
        </details>
      );

    case "skill_start":
      return (
        <div className={styles.skillStart}>
          Starting skill: /{message.skill}
        </div>
      );

    case "skill_complete":
      return (
        <div className={styles.skillComplete}>
          Completed: /{message.skill}
        </div>
      );

    case "result":
      return (
        <div
          className={
            message.status === "error"
              ? styles.resultError
              : styles.resultSuccess
          }
        >
          <strong>
            {message.status === "error" ? "Error" : "Execution complete"}
          </strong>
          {message.cost != null && (
            <span className={styles.resultCost}>
              Cost: ${message.cost.toFixed(4)}
            </span>
          )}
        </div>
      );

    case "error":
      return <div className={styles.errorBlock}>{message.text}</div>;

    default:
      return null;
  }
}
