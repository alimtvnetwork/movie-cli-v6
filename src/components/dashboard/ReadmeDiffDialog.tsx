import { useMemo, useRef, useState } from "react";
import { diffLines, type Change } from "diff";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Upload, FileText, Plus, Minus, Equal } from "lucide-react";
import { toast } from "sonner";

interface ReadmeDiffDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  newContent: string;
  onConfirm: () => void;
}

interface DiffStats {
  added: number;
  removed: number;
  unchanged: number;
}

function computeStats(changes: Change[]): DiffStats {
  let added = 0;
  let removed = 0;
  let unchanged = 0;
  for (const part of changes) {
    const lines = part.value.split("\n").filter((l, i, arr) => !(i === arr.length - 1 && l === "")).length;
    if (part.added) added += lines;
    else if (part.removed) removed += lines;
    else unchanged += lines;
  }
  return { added, removed, unchanged };
}

export function ReadmeDiffDialog({ open, onOpenChange, newContent, onConfirm }: ReadmeDiffDialogProps) {
  const [oldContent, setOldContent] = useState<string>("");
  const [hasLoadedOld, setHasLoadedOld] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const changes = useMemo(() => diffLines(oldContent, newContent), [oldContent, newContent]);
  const stats = useMemo(() => computeStats(changes), [changes]);

  const handleLoadExisting = () => {
    fileInputRef.current?.click();
  };

  const handleFileChosen = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const text = await file.text();
    setOldContent(text);
    setHasLoadedOld(true);
    toast.success("Loaded existing README.md", {
      description: `${text.split("\n").length} lines compared`,
    });
    e.target.value = "";
  };

  const handleConfirm = () => {
    onConfirm();
    onOpenChange(false);
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="max-w-4xl">
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5 text-primary" />
            Confirm README.md changes
          </AlertDialogTitle>
          <AlertDialogDescription>
            Review the exact diff between {hasLoadedOld ? "your existing" : "an empty"} README.md and the new content
            below before saving.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="outline" className="gap-1 border-emerald-500/40 text-emerald-600 dark:text-emerald-400">
            <Plus className="h-3 w-3" />
            {stats.added} added
          </Badge>
          <Badge variant="outline" className="gap-1 border-destructive/40 text-destructive">
            <Minus className="h-3 w-3" />
            {stats.removed} removed
          </Badge>
          <Badge variant="outline" className="gap-1 text-muted-foreground">
            <Equal className="h-3 w-3" />
            {stats.unchanged} unchanged
          </Badge>
          <div className="ml-auto">
            <input
              ref={fileInputRef}
              type="file"
              accept=".md,text/markdown,text/plain"
              className="hidden"
              onChange={handleFileChosen}
            />
            <Button variant="outline" size="sm" className="gap-1.5" onClick={handleLoadExisting}>
              <Upload className="h-3.5 w-3.5" />
              {hasLoadedOld ? "Replace existing" : "Load existing README.md"}
            </Button>
          </div>
        </div>

        <Tabs defaultValue="diff" className="w-full">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="diff">Diff</TabsTrigger>
            <TabsTrigger value="new">New content</TabsTrigger>
            <TabsTrigger value="old">Existing</TabsTrigger>
          </TabsList>

          <TabsContent value="diff" className="mt-3">
            <div className="max-h-[60vh] overflow-auto rounded-md border bg-muted/30 font-mono text-xs">
              {changes.length === 0 || (changes.length === 1 && !changes[0].added && !changes[0].removed) ? (
                <div className="p-4 text-muted-foreground">No differences detected.</div>
              ) : (
                changes.map((part, idx) => {
                  const lines = part.value.split("\n");
                  if (lines[lines.length - 1] === "") lines.pop();
                  const bgClass = part.added
                    ? "bg-emerald-500/10 text-emerald-700 dark:text-emerald-300"
                    : part.removed
                      ? "bg-destructive/10 text-destructive"
                      : "text-muted-foreground";
                  const prefix = part.added ? "+" : part.removed ? "-" : " ";
                  return (
                    <div key={idx} className={bgClass}>
                      {lines.map((line, i) => (
                        <div key={i} className="flex">
                          <span className="select-none px-2 opacity-60">{prefix}</span>
                          <span className="whitespace-pre-wrap break-words">{line || "\u00A0"}</span>
                        </div>
                      ))}
                    </div>
                  );
                })
              )}
            </div>
          </TabsContent>

          <TabsContent value="new" className="mt-3">
            <pre className="max-h-[60vh] overflow-auto rounded-md border bg-muted/30 p-4 font-mono text-xs whitespace-pre-wrap break-words">
              {newContent}
            </pre>
          </TabsContent>

          <TabsContent value="old" className="mt-3">
            <pre className="max-h-[60vh] overflow-auto rounded-md border bg-muted/30 p-4 font-mono text-xs whitespace-pre-wrap break-words">
              {oldContent || "(no existing README.md loaded — diff is against an empty file)"}
            </pre>
          </TabsContent>
        </Tabs>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={handleConfirm}>Confirm & save README.md</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
