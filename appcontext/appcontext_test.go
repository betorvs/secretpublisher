package appcontext

import "testing"

func TestContext_Add(t *testing.T) {
	type fields struct {
		components map[string]Component
	}
	type args struct {
		componentName string
		component     Component
	}
	components := make(map[string]Component)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Add complete repository",
			fields: fields{components: components},
			args: args{
				componentName: Repository,
				component:     ApplicationContext{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := CreateApplicationContext()
			context.Add(tt.args.componentName, tt.args.component)
			if context.Count() == 0 {
				t.Error("Component not added")
			}
			Repository :=
				context.Get(Repository)
			if Repository == nil {
				t.Error("Component not found")
			}
			context.Delete("Repository")
			Repository =
				context.Get("Repository")
			if Repository != nil {
				t.Error("Component not deleted")
			}
		})
	}

}
