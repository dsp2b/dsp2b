package command

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"os"
	"strings"
)

func initCmd(cmd *cobra.Command, args []string) error {
	// 检查是否已经初始化
	if _, err := os.Stat("dspb.json"); err == nil {
		logger.Default().Error("仓库已经初始化")
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	// 初始化仓库
	// 要求用户输入仓库信息
	p := tea.NewProgram(newInitModel())
	_, err := p.Run()
	if err != nil {
		return err
	}
	return nil
}

type initModel struct {
	ti   textinput.Model
	repo *Repository
}

func newInitModel() tea.Model {
	ti := textinput.New()
	ti.Focus()
	return &initModel{
		ti:   ti,
		repo: &Repository{},
	}
}
func (m initModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// 读取仓库信息
			ids := strings.Split(m.ti.Value(), "/")
			id := ids[len(ids)-1]
			objectId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				logger.Default().Error("请输入正确的蓝图集链接")
				return m, tea.Quit
			}
			// 加载蓝图集
			if err := cloneRepo(objectId); err != nil {
				logger.Default().Error("加载蓝图集失败", zap.Error(err))
				return m, tea.Quit
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return m, textinput.Blink
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m *initModel) View() string {
	return fmt.Sprintf(
		"请输入蓝图集链接（例如：https://www.dsp2b.com/zh-CN/collection/65ae4520a7dfce0dabd78905）：\n%s\n%s",
		m.ti.View(),
		"(按esc退出)",
	) + "\n"
}
