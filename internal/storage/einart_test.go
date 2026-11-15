package storage
import("testing";"time";"os";"path/filepath")
func TestNOKAndStableIDs(t *testing.T){
 store,err:=InitializeJsonStore(SystemConfig{StorageURL:t.TempDir()});if err!=nil{t.Fatal(err)}
 if err=store.UpdateCurrency("nok");err!=nil{t.Fatal(err)}
 row:=Expense{ID:"source-row-1",Name:"Lunch",Category:"Food",Amount:-12.34,Currency:"nok",Date:time.Date(2026,11,15,0,0,0,0,time.UTC)}
 if err=store.AddExpense(row);err!=nil{t.Fatal(err)}
 if err=store.AddExpense(row);err!=nil{t.Fatal(err)}
 rows,_:=store.GetAllExpenses();if len(rows)!=1{t.Fatal("repeat created a duplicate")}
 row.Amount=-13;if store.AddExpense(row)==nil{t.Fatal("changed source ID accepted")}
}
func TestAtomicSnapshot(t *testing.T){
 path:=filepath.Join(t.TempDir(),"state.json")
 if err:=atomicWrite(path,[]byte(`{"ok":true}`));err!=nil{t.Fatal(err)}
 if err:=atomicWrite(path,[]byte(`{"next":true}`));err!=nil{t.Fatal(err)}
 b,_:=os.ReadFile(path);if string(b)!=`{"next":true}`{t.Fatal(string(b))}
 info,_:=os.Stat(path);if info.Mode().Perm()!=0600{t.Fatal(info.Mode())}
}
