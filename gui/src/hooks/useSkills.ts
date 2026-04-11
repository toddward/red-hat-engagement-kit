import { useState, useEffect, useCallback } from "react";
import type { SkillSummary, SkillDetail } from "../types";

export function useSkills() {
  const [skills, setSkills] = useState<SkillSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/skills")
      .then((res) => res.json())
      .then((data) => setSkills(data.skills))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const getSkillDetail = useCallback(
    async (name: string): Promise<SkillDetail | null> => {
      try {
        const res = await fetch(`/api/skills/${name}`);
        if (!res.ok) return null;
        return await res.json();
      } catch {
        return null;
      }
    },
    []
  );

  return { skills, loading, error, getSkillDetail };
}
