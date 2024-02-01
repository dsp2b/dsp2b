package command

import "github.com/spf13/cobra"

func AddCommand(cmd *cobra.Command) {
	cmd.AddCommand([]*cobra.Command{{
		Use:   "init",
		Short: "克隆戴森球蓝图仓库",
		RunE:  initCmd,
	}, {
		Use:   "push",
		Short: "发布戴森球蓝图集",
		RunE:  push,
	}, {
		Use:   "pull",
		Short: "拉取戴森球蓝图集",
		RunE:  pull,
	}, {
		Use:   "clone [collection_id]",
		Short: "克隆戴森球蓝图集",
		RunE:  cloneCmd,
	}, {
		Use:   "status",
		Short: "查看戴森球蓝图集状态",
		RunE:  statusCmd,
	}, {
		Use:   "rename",
		Short: "重命名所有蓝图",
		RunE:  renameCmd,
	}}...)
}
