package markitdown

type Options struct {
	DataDir string `json:"data_dir"`
}

func (p *Options) toCommandArgs() (ret []string, cleanups []func(), err error) {
	var args []string
	var cleanupFuncs []func()

	defer func() {
		if err != nil {
			for i := 0; i < len(cleanupFuncs); i++ {
				cleanupFuncs[i]()
			}
		}
	}()
	return args, cleanupFuncs, nil
}
