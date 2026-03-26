import { Modal } from "@arco-design/web-vue";

type ConfirmOptions = {
  title: string;
  content: string;
  okText?: string;
  cancelText?: string;
  okButtonProps?: Record<string, unknown>;
};

const askConfirmation = (options: ConfirmOptions) =>
  new Promise<boolean>(resolve => {
    Modal.confirm({
      title: options.title,
      content: options.content,
      okText: options.okText || "确认",
      cancelText: options.cancelText || "取消",
      okButtonProps: options.okButtonProps,
      onOk: () => resolve(true),
      onCancel: () => resolve(false)
    });
  });

export function useConfirmAction() {
  const runConfirmed = async (options: ConfirmOptions, action: () => Promise<void> | void) => {
    const confirmed = await askConfirmation(options);
    if (!confirmed) {
      return false;
    }
    await action();
    return true;
  };

  const confirmSave = (action: () => Promise<void> | void, target = "当前内容") =>
    runConfirmed(
      {
        title: "确认保存",
        content: `确认保存${target}吗？`
      },
      action
    );

  const confirmDelete = (action: () => Promise<void> | void, target = "当前数据") =>
    runConfirmed(
      {
        title: "确认删除",
        content: `删除后将无法恢复，确认删除${target}吗？`,
        okButtonProps: { status: "danger" }
      },
      action
    );

  const confirmBatchDelete = (action: () => Promise<void> | void, target = "选中数据") =>
    runConfirmed(
      {
        title: "确认批量删除",
        content: `删除后将无法恢复，确认批量删除${target}吗？`,
        okButtonProps: { status: "danger" }
      },
      action
    );

  const confirmBatchHide = (action: () => Promise<void> | void, target = "选中文章") =>
    runConfirmed(
      {
        title: "确认批量隐藏",
        content: `确认将${target}批量设置为隐藏吗？`
      },
      action
    );

  return {
    runConfirmed,
    confirmSave,
    confirmDelete,
    confirmBatchDelete,
    confirmBatchHide
  };
}
