import { useState, useEffect, useCallback } from "react";
import type { Engagement } from "../types";

export function useEngagements() {
  const [engagements, setEngagements] = useState<Engagement[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(() => {
    setLoading(true);
    fetch("/api/engagements")
      .then((res) => res.json())
      .then((data) => setEngagements(data.engagements))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const getContext = useCallback(
    async (slug: string): Promise<string | null> => {
      try {
        const res = await fetch(`/api/engagements/${slug}/context`);
        if (!res.ok) return null;
        const data = await res.json();
        return data.content;
      } catch {
        return null;
      }
    },
    []
  );

  return { engagements, loading, error, getContext, refresh };
}
