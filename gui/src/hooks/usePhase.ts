import { useState, useEffect, useCallback } from "react";
import type { PhaseInfo } from "../types";

export function usePhase(slug: string | null) {
  const [phase, setPhase] = useState<PhaseInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(() => {
    if (!slug) {
      setPhase(null);
      return;
    }
    setLoading(true);
    fetch(`/api/engagements/${slug}/phase`)
      .then((res) => res.json())
      .then((data) => setPhase(data))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [slug]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return { phase, loading, error, refresh };
}
