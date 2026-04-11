import type { ReactNode } from "react";
import styles from "./Layout.module.css";

interface LayoutProps {
  sidebar: ReactNode;
  main: ReactNode;
}

export default function Layout({ sidebar, main }: LayoutProps) {
  return (
    <div className={styles.container}>
      <aside className={styles.sidebar}>
        <div className={styles.sidebarHeader}>
          <h1 className={styles.sidebarTitle}>Red Hat Engagement Kit</h1>
        </div>
        <div className={styles.sidebarContent}>{sidebar}</div>
      </aside>
      <main className={styles.main}>{main}</main>
    </div>
  );
}
