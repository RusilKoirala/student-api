import { useEffect, useState } from "react";
import { Plus, Trash2, Users, Edit2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
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
import { getTeachers, createTeacher, deleteTeacher, updateTeacher } from "@/lib/api";

const EMPTY = { name: "", email: "", subject: "" };

export default function Teachers() {
  const [teachers, setTeachers] = useState([]);
  const [loading, setLoading]   = useState(true);
  const [open, setOpen]         = useState(false);
  const [form, setForm]         = useState(EMPTY);
  const [editingId, setEditingId] = useState(null);
  const [saving, setSaving]     = useState(false);

  const load = () =>
    getTeachers()
      .then(setTeachers)
      .finally(() => setLoading(false));

  useEffect(() => { load(); }, []);

  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  async function handleSave(e) {
    e.preventDefault();
    if (!form.name || !form.email || !form.subject) {
      toast.error("All fields are required");
      return;
    }
    setSaving(true);
    try {
      if (editingId) {
        await updateTeacher(editingId, form);
        toast.success("Teacher updated");
      } else {
        await createTeacher(form);
        toast.success("Teacher added");
      }
      setOpen(false);
      setForm(EMPTY);
      setEditingId(null);
      load();
    } catch {
      toast.error(editingId ? "Failed to update teacher" : "Failed to add teacher");
    } finally {
      setSaving(false);
    }
  }

  function handleEdit(teacher) {
    setEditingId(teacher.id);
    setForm({
      name: teacher.name,
      email: teacher.email,
      subject: teacher.subject,
    });
    setOpen(true);
  }

  async function handleDelete(id, name) {
    if (!confirm(`Delete ${name}?`)) return;
    try {
      await deleteTeacher(id);
      toast.success("Teacher deleted");
      load();
    } catch {
      toast.error("Failed to delete teacher");
    }
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">Teachers</h2>
          <p className="text-sm text-muted-foreground mt-0.5">
            {teachers.length} teacher{teachers.length !== 1 ? "s" : ""} total
          </p>
        </div>
        <Button onClick={() => { setEditingId(null); setForm(EMPTY); setOpen(true); }}>
          <Plus className="size-4" />
          Add Teacher
        </Button>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Subject</TableHead>
              <TableHead className="w-24" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={4} className="h-32 text-center text-muted-foreground">
                  Loading...
                </TableCell>
              </TableRow>
            ) : teachers.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="h-32 text-center">
                  <div className="flex flex-col items-center gap-2 text-muted-foreground">
                    <Users className="size-8 opacity-40" />
                    <p className="text-sm">No teachers yet</p>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              teachers.map((t) => (
                <TableRow key={t.id}>
                  <TableCell className="font-medium">{t.name}</TableCell>
                  <TableCell className="text-muted-foreground">{t.email}</TableCell>
                  <TableCell>
                    <Badge variant="secondary">{t.subject}</Badge>
                  </TableCell>
                  <TableCell className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      onClick={() => handleEdit(t)}
                    >
                      <Edit2 className="size-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      className="text-destructive hover:text-destructive"
                      onClick={() => handleDelete(t.id, t.name)}
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
            <DialogTitle>{editingId ? "Edit Teacher" : "Add Teacher"}</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="t-name">Full Name</Label>
              <Input id="t-name" placeholder="Jane Smith" value={form.name} onChange={set("name")} />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="t-email">Email</Label>
              <Input id="t-email" type="email" placeholder="jane@school.edu" value={form.email} onChange={set("email")} />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="t-subject">Subject</Label>
              <Input id="t-subject" placeholder="Mathematics" value={form.subject} onChange={set("subject")} />
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => { setOpen(false); setEditingId(null); setForm(EMPTY); }}>
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                {saving ? "Saving…" : (editingId ? "Save Changes" : "Add Teacher")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
