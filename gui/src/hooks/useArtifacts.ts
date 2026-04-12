import { useState, useEffect, useCallback } from "react";
import type { ArtifactNode } from "../types";

export function useArtifacts(slug: string | null) {
  const [tree, setTree] = useState<ArtifactNode[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(() => {
    if (!slug) {
      setTree(null);
      return;
    }
    setLoading(true);
    fetch(`/api/engagements/${slug}/artifacts`)
      .then((res) => res.json())
      .then((data) => setTree(data.tree))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [slug]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return { tree, loading, error, refresh };
}

export function useArtifactContent(slug: string | null, filePath: string | null) {
  const [content, setContent] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!slug || !filePath) {
      setContent(null);
      return;
    }
    setLoading(true);
    fetch(`/api/engagements/${slug}/artifacts/${filePath}`)
      .then((res) => res.json())
      .then((data) => setContent(data.content))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [slug, filePath]);

  return { content, loading, error };
}
