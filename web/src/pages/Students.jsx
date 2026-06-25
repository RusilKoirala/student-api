import { useEffect, useState } from "react";
import { Plus, Trash2, GraduationCap } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
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
import { getStudents, createStudent, deleteStudent, getClasses } from "@/lib/api";

const EMPTY = { name: "", email: "", age: "", classId: "" };

export default function Students() {
  const [students, setStudents] = useState([]);
  const [classes,  setClasses]  = useState([]);
  const [loading, setLoading]   = useState(true);
  const [open, setOpen]         = useState(false);
  const [form, setForm]         = useState(EMPTY);
  const [saving, setSaving]     = useState(false);

  const classMap = Object.fromEntries(classes.map((c) => [c.id, c.name]));

  const load = () =>
    Promise.all([getStudents(), getClasses()])
      .then(([s, c]) => { setStudents(s); setClasses(c); })
      .finally(() => setLoading(false));

  useEffect(() => { load(); }, []);

  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  async function handleCreate(e) {
    e.preventDefault();
    if (!form.name || !form.email || !form.age) {
      toast.error("Name, email and age are required");
      return;
    }
    setSaving(true);
    try {
      await createStudent({
        name:    form.name,
        email:   form.email,
        age:     Number(form.age),
        classId: Number(form.classId) || 0,
      });
      toast.success("Student enrolled");
      setOpen(false);
      setForm(EMPTY);
      load();
    } catch {
      toast.error("Failed to enrol student");
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id, name) {
    if (!confirm(`Remove ${name} from the school?`)) return;
    try {
      await deleteStudent(id);
      toast.success("Student removed");
      load();
    } catch {
      toast.error("Failed to remove student");
    }
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">Students</h2>
          <p className="text-sm text-muted-foreground mt-0.5">
            {students.length} student{students.length !== 1 ? "s" : ""} enrolled
          </p>
        </div>
        <Button onClick={() => setOpen(true)}>
          <Plus className="size-4" />
          Enrol Student
        </Button>
      </div>

      {/* Table */}
      <div className="rounded-lg border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Age</TableHead>
              <TableHead>Class</TableHead>
              <TableHead className="w-16" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                  Loading...
                </TableCell>
              </TableRow>
            ) : students.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="h-32 text-center">
                  <div className="flex flex-col items-center gap-2 text-muted-foreground">
                    <GraduationCap className="size-8 opacity-40" />
                    <p className="text-sm">No students enrolled yet</p>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              students.map((s) => (
                <TableRow key={s.id}>
                  <TableCell className="font-medium">{s.name}</TableCell>
                  <TableCell className="text-muted-foreground">{s.email}</TableCell>
                  <TableCell>{s.age}</TableCell>
                  <TableCell>
                    {s.classId && classMap[s.classId] ? (
                      <Badge variant="secondary">{classMap[s.classId]}</Badge>
                    ) : (
                      <span className="text-sm italic text-muted-foreground">None</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      className="text-destructive hover:text-destructive"
                      onClick={() => handleDelete(s.id, s.name)}
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

      {/* Enrol Dialog */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Enrol Student</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleCreate} className="space-y-4">
            <div className="space-y-1.5">
              <Label htmlFor="s-name">Full Name</Label>
              <Input id="s-name" placeholder="John Doe" value={form.name} onChange={set("name")} />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="s-email">Email</Label>
              <Input id="s-email" type="email" placeholder="john@school.edu" value={form.email} onChange={set("email")} />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="s-age">Age</Label>
              <Input id="s-age" type="number" min="3" max="30" placeholder="15" value={form.age} onChange={set("age")} />
            </div>
            <div className="space-y-1.5">
              <Label>Class <span className="text-muted-foreground">(optional)</span></Label>
              <Select
                value={form.classId}
                onValueChange={(v) => setForm((f) => ({ ...f, classId: v }))}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select a class…" />
                </SelectTrigger>
                <SelectContent>
                  {classes.map((c) => (
                    <SelectItem key={c.id} value={String(c.id)}>
                      {c.name}
                      {c.teacherName ? ` — ${c.teacherName}` : ""}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                {saving ? "Saving…" : "Enrol"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
