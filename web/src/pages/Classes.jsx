import { useEffect, useState } from "react";
import { Plus, Trash2, BookOpen, Edit2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getClasses, createClass, deleteClass, updateClass, getTeachers } from "@/lib/api";

const EMPTY = { name: "", teacherId: "" };

export default function Classes() {
  const [classes,  setClasses]  = useState([]);
  const [teachers, setTeachers] = useState([]);
  const [loading, setLoading]   = useState(true);
  const [open, setOpen]         = useState(false);
  const [form, setForm]         = useState(EMPTY);
  const [editingId, setEditingId] = useState(null);
  const [saving, setSaving]     = useState(false);

  const load = () =>
    Promise.all([getClasses(), getTeachers()])
      .then(([c, t]) => { setClasses(c); setTeachers(t); })
      .finally(() => setLoading(false));

  useEffect(() => { load(); }, []);

  async function handleSave(e) {
    e.preventDefault();
    if (!form.name) { toast.error("Class name is required"); return; }
    setSaving(true);
    try {
      if (editingId) {
        await updateClass(editingId, { name: form.name, teacherId: Number(form.teacherId) || 0 });
        toast.success("Class updated");
      } else {
        await createClass({ name: form.name, teacherId: Number(form.teacherId) || 0 });
        toast.success("Class created");
      }
      setOpen(false);
      setForm(EMPTY);
      setEditingId(null);
      load();
    } catch {
      toast.error(editingId ? "Failed to update class" : "Failed to create class");
    } finally {
      setSaving(false);
    }
  }

  function handleEdit(cls) {
    setEditingId(cls.id);
    setForm({
      name: cls.name,
      teacherId: cls.teacherId ? String(cls.teacherId) : "",
    });
    setOpen(true);
  }

  async function handleDelete(id, name) {
    if (!confirm(`Delete class "${name}"?`)) return;
    try {
      await deleteClass(id);
      toast.success("Class deleted");
      load();
    } catch {
      toast.error("Failed to delete class");
    }
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">Classes</h2>
          <p className="text-sm text-muted-foreground mt-0.5">
            {classes.length} class{classes.length !== 1 ? "es" : ""} total
          </p>
        </div>
        <Button onClick={() => { setEditingId(null); setForm(EMPTY); setOpen(true); }}>
          <Plus className="size-4" />
          Add Class
        </Button>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Class Name</TableHead>
              <TableHead>Teacher</TableHead>
              <TableHead className="w-24" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={3} className="h-32 text-center text-muted-foreground">
                  Loading...
                </TableCell>
              </TableRow>
            ) : classes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={3} className="h-32 text-center">
                  <div className="flex flex-col items-center gap-2 text-muted-foreground">
                    <BookOpen className="size-8 opacity-40" />
                    <p className="text-sm">No classes yet</p>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              classes.map((c) => (
                <TableRow key={c.id}>
                  <TableCell className="font-medium">{c.name}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {c.teacherName || <span className="italic opacity-50">Unassigned</span>}
                  </TableCell>
                  <TableCell className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      onClick={() => handleEdit(c)}
                    >
                      <Edit2 className="size-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      className="text-destructive hover:text-destructive"
                      onClick={() => handleDelete(c.id, c.name)}
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Add/Edit Dialog */}
      <Dialog open={open} onOpenChange={(o) => { if (!o) { setEditingId(null); setForm(EMPTY); } setOpen(o); }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editingId ? "Edit Class" : "Add Class"}</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="c-name">Class Name</Label>
              <Input
                id="c-name"
                placeholder="e.g. Grade 10A"
                value={form.name}
                onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              />
            </div>
            <div className="space-y-1.5">
              <Label>Assign Teacher <span className="text-muted-foreground">(optional)</span></Label>
              <Select
                value={form.teacherId}
                onValueChange={(v) => setForm((f) => ({ ...f, teacherId: v }))}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a teacher…" />
                </SelectTrigger>
                <SelectContent>
                  {teachers.map((t) => (
                    <SelectItem key={t.id} value={String(t.id)}>
                      {t.name} — {t.subject}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setOpen(false); setEditingId(null); setForm(EMPTY); }}>
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                {saving ? "Saving…" : (editingId ? "Save Changes" : "Add Class")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
