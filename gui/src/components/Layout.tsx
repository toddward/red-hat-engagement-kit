import type { ReactNode } from "react";
import redhatLogo from "../assets/redhat-logo.svg";
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
          <img src={redhatLogo} alt="Red Hat" className={styles.logo} />
          <h1 className={styles.sidebarTitle}>Engagement Kit</h1>
        </div>
        <div className={styles.sidebarContent}>{sidebar}</div>
      </aside>
      <main className={styles.main}>{main}</main>
    </div>
  );
}
