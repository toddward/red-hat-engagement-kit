import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import styles from "./SkillDetail.module.css";

interface SkillDetailProps {
  name: string;
  content: string;
}

export default function SkillDetail({ name, content }: SkillDetailProps) {
  return (
    <div className={styles.container}>
      <div className={styles.label}>Skill Definition</div>
      <h2 className={styles.title}>/{name}</h2>
      <div className={styles.markdownBody}>
        <Markdown remarkPlugins={[remarkGfm]}>{content}</Markdown>
      </div>
    </div>
  );
}
