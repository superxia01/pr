-- 充值订单表
-- 用于线下充值流程：用户提交订单 → 超管审核 → 充值入账

CREATE TABLE IF NOT EXISTS recharge_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    account_id UUID NOT NULL REFERENCES credit_accounts(id),
    amount INT NOT NULL CHECK (amount > 0),
    payment_method VARCHAR(20) NOT NULL CHECK (payment_method IN ('alipay', 'wechat', 'bank')),
    payment_proof VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'completed')),
    rejection_note TEXT,
    audited_by VARCHAR(255) REFERENCES users(id),
    audited_at TIMESTAMP,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_recharge_orders_user ON recharge_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_account ON recharge_orders(account_id);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_status ON recharge_orders(status);
CREATE INDEX IF NOT EXISTS idx_recharge_orders_audited_by ON recharge_orders(audited_by);

-- 注释
COMMENT ON TABLE recharge_orders IS '充值订单表（线下充值流程）';
COMMENT ON COLUMN recharge_orders.user_id IS '提交订单的用户ID';
COMMENT ON COLUMN recharge_orders.account_id IS '要充值的积分账户ID';
COMMENT ON COLUMN recharge_orders.amount IS '充值金额（积分）';
COMMENT ON COLUMN recharge_orders.payment_method IS '支付方式：alipay-支付宝, wechat-微信, bank-银行转账';
COMMENT ON COLUMN recharge_orders.payment_proof IS '支付凭证URL（用户上传的转账截图等）';
COMMENT ON COLUMN recharge_orders.status IS '订单状态：pending-待审核, approved-已通过, rejected-已拒绝, completed-已完成';
COMMENT ON COLUMN recharge_orders.rejection_note IS '拒绝原因（当status=rejected时）';
COMMENT ON COLUMN recharge_orders.audited_by IS '审核人ID（超管）';
COMMENT ON COLUMN recharge_orders.audited_at IS '审核时间';
COMMENT ON COLUMN recharge_orders.processed_at IS '完成时间（充值入账时间）';
