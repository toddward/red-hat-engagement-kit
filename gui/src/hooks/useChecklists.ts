import { useState, useEffect, useCallback, useMemo } from "react";
import type { Checklist } from "../types";

export function useChecklists(slug: string | null) {
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(() => {
    if (!slug) {
      setChecklists([]);
      return;
    }
    setLoading(true);
    fetch(`/api/engagements/${slug}/checklists`)
      .then((res) => res.json())
      .then((data) => setChecklists(data.checklists))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [slug]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const instantiate = useCallback(async () => {
    if (!slug) return;
    await fetch(`/api/engagements/${slug}/checklists/instantiate`, {
      method: "POST",
    });
    refresh();
  }, [slug, refresh]);

  const toggleItem = useCallback(
    async (fileName: string, lineNumber: number, checked: boolean) => {
      if (!slug) return;

      // Optimistic update
      setChecklists((prev) =>
        prev.map((cl) => {
          if (cl.fileName !== fileName) return cl;
          let total = 0;
          let done = 0;
          const sections = cl.sections.map((sec) => ({
            ...sec,
            items: sec.items.map((item) => {
              const updated =
                item.line === lineNumber ? { ...item, checked } : item;
              total++;
              if (updated.checked) done++;
              return updated;
            }),
          }));
          return {
            ...cl,
            sections,
            completionPercent: total > 0 ? Math.round((done / total) * 100) : 0,
          };
        })
      );

      try {
        const res = await fetch(
          `/api/engagements/${slug}/checklists/toggle`,
          {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ fileName, lineNumber, checked }),
          }
        );
        if (!res.ok) {
          refresh(); // Revert on failure
        }
      } catch {
        refresh(); // Revert on error
      }
    },
    [slug, refresh]
  );

  const summary = useMemo(() => {
    let total = 0;
    let checked = 0;
    for (const cl of checklists) {
      for (const sec of cl.sections) {
        for (const item of sec.items) {
          total++;
          if (item.checked) checked++;
        }
      }
    }
    return {
      total,
      checked,
      percent: total > 0 ? Math.round((checked / total) * 100) : 0,
    };
  }, [checklists]);

  return { checklists, loading, error, refresh, instantiate, toggleItem, summary };
}
