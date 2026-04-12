import styles from "./TypingIndicator.module.css";

export default function TypingIndicator() {
  return (
    <div className={styles.container}>
      <div className={styles.dot} />
      <div className={styles.dot} />
      <div className={styles.dot} />
    </div>
  );
}
