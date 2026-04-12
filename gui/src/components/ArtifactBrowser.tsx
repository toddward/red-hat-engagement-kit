import { useState } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useArtifacts, useArtifactContent } from "../hooks/useArtifacts";
import type { ArtifactNode } from "../types";
import styles from "./ArtifactBrowser.module.css";
import detailStyles from "./SkillDetail.module.css";

interface ArtifactBrowserProps {
  slug: string;
}

export default function ArtifactBrowser({ slug }: ArtifactBrowserProps) {
  const { tree, loading } = useArtifacts(slug);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const { content } = useArtifactContent(slug, selectedFile);

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.emptyContent}>Loading artifacts...</div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.treePanel}>
        <div className={styles.treePanelHeader}>Engagement Files</div>
        {tree && tree.length > 0 ? (
          tree.map((node) => (
            <FileTreeNode
              key={node.path}
              node={node}
              selectedFile={selectedFile}
              onSelect={setSelectedFile}
              depth={0}
            />
          ))
        ) : (
          <div style={{ fontSize: 12, color: "var(--rh-gray-500)" }}>
            No artifacts found.
          </div>
        )}
      </div>

      <div className={styles.contentPanel}>
        {content ? (
          <>
            <div className={styles.contentHeader}>Artifact</div>
            <div className={styles.contentTitle}>{selectedFile}</div>
            <div className={detailStyles.markdownBody}>
              <Markdown remarkPlugins={[remarkGfm]}>{content}</Markdown>
            </div>
          </>
        ) : (
          <div className={styles.emptyContent}>
            Select a file to view its contents.
          </div>
        )}
      </div>
    </div>
  );
}

function FileTreeNode({
  node,
  selectedFile,
  onSelect,
  depth,
}: {
  node: ArtifactNode;
  selectedFile: string | null;
  onSelect: (path: string) => void;
  depth: number;
}) {
  const [expanded, setExpanded] = useState(depth < 1);
  const isDir = node.type === "directory";
  const isActive = selectedFile === node.path;

  const handleClick = () => {
    if (isDir) {
      setExpanded((prev) => !prev);
    } else {
      onSelect(node.path);
    }
  };

  return (
    <div className={styles.treeNode}>
      <div
        className={`${styles.treeRow} ${isActive ? styles.treeRowActive : ""}`}
        onClick={handleClick}
      >
        <span className={styles.treeIcon}>
          {isDir ? (expanded ? "\u25BE" : "\u25B8") : "\u2022"}
        </span>
        <span className={styles.treeName}>{node.name}</span>
      </div>
      {isDir && expanded && node.children && node.children.length > 0 && (
        <div className={styles.treeChildren}>
          {node.children.map((child) => (
            <FileTreeNode
              key={child.path}
              node={child}
              selectedFile={selectedFile}
              onSelect={onSelect}
              depth={depth + 1}
            />
          ))}
        </div>
      )}
    </div>
  );
}
