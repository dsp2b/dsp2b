package command

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/codfrm/cago/pkg/logger"
	"github.com/codfrm/cago/pkg/utils/httputils"
	api "github.com/dsp2b/dsp2b-go/internal/api/blueprint"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

const BaseURL = "https://www.dsp2b.com/zh-CN"

var configDir = "~/.dsp2b"

// Abs 转绝对路径,处理了"~"
func Abs(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Abs(filepath.Join(home, p[1:]))
	}
	return filepath.Abs(p)
}

func init() {
	// 判断是否为windows
	if runtime.GOOS == "windows" {
		configDir, _ = os.UserHomeDir()
	}
	configDir, _ = Abs(configDir)
}

type ApiClient struct {
	cookie string
}

func NewApiClient() (*ApiClient, error) {
	cookie, err := os.ReadFile(filepath.Join(configDir, ".cookie"))
	if err != nil {
		logger.Default().Error("读取cookie失败", zap.Error(err))
		return nil, err
	}
	return &ApiClient{
		cookie: string(cookie),
	}, nil
}

func (a *ApiClient) request(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, BaseURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", a.cookie)
	return req, nil
}

func (a *ApiClient) GET(ctx context.Context, url string) (*http.Response, error) {
	req, err := a.request(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ApiClient) POST(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := a.request(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CollectionDetail(ctx context.Context, id primitive.ObjectID, page int) (*CollectionDetailResponse, error) {
	resp, err := http.Get(BaseURL + "/collection/" + id.Hex() +
		"?page=" + strconv.Itoa(page) + "&_data=routes%2F%24lng.collection_.%24id&root=false")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Ctx(ctx).Info("请求失败", zap.String("body", string(data)))
		return nil, errors.New("请求失败")
	}
	collection := &CollectionDetailResponse{}
	if err := json.Unmarshal(data, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

func ReadAllBlueprint(ctx context.Context, id primitive.ObjectID) (*CollectionDetailResponse, []*BlueprintItem, error) {
	var items []*BlueprintItem
	collection, err := CollectionDetail(ctx, id, 1)
	if err != nil {
		return nil, nil, err
	}
	items = append(items, collection.List...)
	if len(collection.List) == collection.Total {
		return collection, items, nil
	}
	for i := 2; ; i++ {
		collection, err := CollectionDetail(ctx, id, i)
		if err != nil {
			return nil, nil, err
		}
		if len(collection.List) == 0 {
			break
		}
		items = append(items, collection.List...)
	}
	return collection, items, nil
}

func (a *ApiClient) PostBlueprint(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := a.POST(ctx, "/create/blueprint?_data=routes%2F%24lng.create.blueprint.%24%28id%29",
		"application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Ctx(ctx).Info("请求失败", zap.String("body", string(data)))
		return nil, errors.New("请求失败")
	}
	logger.Ctx(ctx).Info("请求成功", zap.String("body", string(data)))
	result := &api.CreateResponse{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ApiClient) ParseBlueprint(ctx context.Context, req *api.ParseRequest) (*api.ParseResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := a.POST(ctx,
		"/create/blueprint?action=parse&_data=routes%2F%24lng.create.blueprint.%24%28id%29",
		"application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Ctx(ctx).Info("请求失败", zap.String("body", string(data)))
		return nil, errors.New("请求失败")
	}
	result := &api.ParseResponse{}
	if err := json.Unmarshal(data, &httputils.JSONResponse{
		Code: 0,
		Msg:  "",
		Data: result,
	}); err != nil {
		return nil, err
	}
	return result, nil
}
